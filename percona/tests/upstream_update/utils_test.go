package percona_tests

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	portRangeStart = 20000 // exporter web interface listening port
	portRangeEnd   = 20100 // exporter web interface listening port

	exporterWaitTimeoutMs = 3000 // time to wait for exporter process start

	updatedExporterFileName = "assets/node_exporter"
	oldExporterFileName     = "assets/node_exporter_percona"
)

func getBool(val *bool) bool {
	return val != nil && *val
}

func launchExporter(fileName string) (cmd *exec.Cmd, port int, collectOutput func() string, _ error) {
	lines, err := os.ReadFile("assets/test.exporter-flags.txt")
	if err != nil {
		return nil, 0, nil, errors.Wrapf(err, "Unable to read exporter args file")
	}

	port = -1
	for i := portRangeStart; i < portRangeEnd; i++ {
		if checkPort(i) {
			port = i
			break
		}
	}

	if port == -1 {
		return nil, 0, nil, errors.Wrapf(err, "Failed to find free port in range [%d..%d]", portRangeStart, portRangeEnd)
	}

	linesStr := string(lines)
	linesStr += fmt.Sprintf("\n--web.listen-address=127.0.0.1:%d", port)

	/*absolutePath, _ := filepath.Abs("custom-queries")
	linesStr += fmt.Sprintf("\n--collect.custom_query.hr.directory=%s/high-resolution", absolutePath)
	linesStr += fmt.Sprintf("\n--collect.custom_query.mr.directory=%s/medium-resolution", absolutePath)
	linesStr += fmt.Sprintf("\n--collect.custom_query.lr.directory=%s/low-resolution", absolutePath)*/

	linesArr := strings.Split(linesStr, "\n")

	cmd = exec.Command(fileName, linesArr...)

	var outBuffer, errorBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errorBuffer

	collectOutput = func() string {
		result := ""
		outStr := outBuffer.String()
		if outStr == "" {
			result = "Process stdOut was empty. "
		} else {
			result = fmt.Sprintf("Process stdOut:\n%s\n", outStr)
		}
		errStr := errorBuffer.String()
		if errStr == "" {
			result += "Process stdErr was empty."
		} else {
			result += fmt.Sprintf("Process stdErr:\n%s\n", errStr)
		}

		return result
	}

	err = cmd.Start()
	if err != nil {
		return nil, 0, nil, errors.Wrapf(err, "Failed to start exporter.%s", collectOutput())
	}

	err = waitForExporter(port)
	if err != nil {
		return nil, 0, nil, errors.Wrapf(err, "Failed to wait for exporter.%s", collectOutput())
	}

	return cmd, port, collectOutput, nil
}

func stopExporter(cmd *exec.Cmd, collectOutput func() string) error {
	err := cmd.Process.Signal(unix.SIGINT)
	if err != nil {
		return errors.Wrapf(err, "Failed to send SIGINT to exporter process.%s\n", collectOutput())
	}

	err = cmd.Wait()
	if err != nil && err.Error() != "signal: interrupt" {
		return errors.Wrapf(err, "Failed to wait for exporter process termination.%s\n", collectOutput())
	}

	return nil
}
func tryGetMetrics(port int) (string, error) {
	return tryGetMetricsFrom(port, "metrics")
}

func tryGetMetricsFrom(port int, endpoint string) (string, error) {
	uri := fmt.Sprintf("http://127.0.0.1:%d/%s", port, endpoint)
	client := new(http.Client)

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Accept-Encoding", "gzip")

	response, err := client.Do(request)

	if err != nil {
		return "", fmt.Errorf("failed to get response from exporters web interface: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get response from exporters web interface: %w", err)
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	enc := response.Header.Get("Content-Encoding")
	switch enc {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return "", fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	rr := buf.String()
	if rr == "" {
		return "", fmt.Errorf("failed to read response")
	}

	err = response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close response: %w", err)
	}

	return rr, nil
}

func checkPort(port int) bool {
	ln, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		return false
	}

	_ = ln.Close()
	return true
}

func waitForExporter(port int) error {
	watchdog := exporterWaitTimeoutMs

	_, e := tryGetMetrics(port)
	for ; e != nil && watchdog > 0; watchdog-- {
		time.Sleep(1 * time.Millisecond)
		_, e = tryGetMetrics(port)
	}

	if watchdog == 0 {
		return fmt.Errorf("failed to wait for exporter (on port %d)", port)
	}

	return nil
}
