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

const (
	binary  = "./node_exporter"
	address = "localhost:19100"
)

func TestFileDescriptorLeak(t *testing.T) {
	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}
	if _, err := procfs.NewStat(); err != nil {
		t.Skipf("proc filesystem is not available, but currently required to read number of open file descriptors: %s", err)
	}

	exporter := exec.Command(binary, "-web.listen-address", address)
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

	if err := runCommandAndTests(exporter, test); err != nil {
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

	exporter := exec.Command(binary, "-web.listen-address", address, "-collector.textfile.directory", dir)
	test := func(_ int) error {
		return queryExporter(address)
	}

	if err := runCommandAndTests(exporter, test); err != nil {
		t.Error(err)
	}
}

func queryExporter(address string) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", address))
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		return fmt.Errorf("want /metrics status code %d, have %d", want, have)
	}
	return nil
}

func runCommandAndTests(cmd *exec.Cmd, fn func(pid int) error) error {
	errc := make(chan error)
	go func() {
		if err := cmd.Run(); err != nil {
			errc <- fmt.Errorf("execution of command failed: %s", err)
		} else {
			errc <- nil
		}
	}()

	// Allow the process to start before running any tests.
	select {
	case err := <-errc:
		return err
	case <-time.After(100 * time.Millisecond):
	}

	go func(pid int) {
		errc <- fn(pid)
	}(cmd.Process.Pid)

	select {
	case err := <-errc:
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		if err != nil {
			return err
		}
	}
	return nil
}
