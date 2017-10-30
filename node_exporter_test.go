package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/procfs"
)

var (
	binary = filepath.Join(os.Getenv("GOPATH"), "bin/node_exporter")
)

const (
	address = "localhost:19100"
)

func TestFileDescriptorLeak(t *testing.T) {
	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}
	if _, err := procfs.NewStat(); err != nil {
		t.Skipf("proc filesystem is not available, but currently required to read number of open file descriptors: %s", err)
	}

	exporter := exec.Command(binary, "--web.listen-address", address)
	test := func(pid int) error {
		if err := queryExporter(address); err != nil {
			return err
		}
		proc, err := procfs.NewProc(pid)
		if err != nil {
			return err
		}
		fdsBefore, err := proc.FileDescriptors()
		if err != nil {
			return err
		}
		for i := 0; i < 5; i++ {
			if err := queryExporter(address); err != nil {
				return err
			}
		}
		fdsAfter, err := proc.FileDescriptors()
		if err != nil {
			return err
		}
		if want, have := len(fdsBefore), len(fdsAfter); want != have {
			return fmt.Errorf("want %d open file descriptors after metrics scrape, have %d", want, have)
		}
		return nil
	}

	if err := runCommandAndTests(exporter, address, test); err != nil {
		t.Error(err)
	}
}

func TestHandlingOfDuplicatedMetrics(t *testing.T) {
	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}

	dir, err := ioutil.TempDir("", "node-exporter")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	content := []byte("dummy_metric 1\n")
	if err := ioutil.WriteFile(filepath.Join(dir, "a.prom"), content, 0600); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "b.prom"), content, 0600); err != nil {
		t.Fatal(err)
	}

	exporter := exec.Command(binary, "--web.listen-address", address, "--collector.textfile.directory", dir)
	test := func(_ int) error {
		return queryExporter(address)
	}

	if err := runCommandAndTests(exporter, address, test); err != nil {
		t.Error(err)
	}
}

func queryExporter(address string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", address))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		return fmt.Errorf("want /metrics status code %d, have %d. Body:\n%s", want, have, b)
	}
	return nil
}

func runCommandAndTests(cmd *exec.Cmd, address string, fn func(pid int) error) error {
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %s", err)
	}
	time.Sleep(50 * time.Millisecond)
	for i := 0; i < 10; i++ {
		if err := queryExporter(address); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
		if cmd.Process == nil || i == 9 {
			return fmt.Errorf("can't start command")
		}
	}

	errc := make(chan error)
	go func(pid int) {
		errc <- fn(pid)
	}(cmd.Process.Pid)

	err := <-errc
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
	return err
}
