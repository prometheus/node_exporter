package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/prometheus/procfs"
)

func TestFileDescriptorLeak(t *testing.T) {
	const (
		binary  = "./node_exporter"
		address = "localhost:9100"
	)

	if _, err := os.Stat(binary); err != nil {
		t.Skipf("node_exporter binary not available, try to run `make build` first: %s", err)
	}
	if _, err := procfs.NewStat(); err != nil {
		t.Skipf("proc filesystem is not available, but currently required to read number of open file descriptors: %s", err)
	}

	errc := make(chan error)
	exporter := exec.Command(binary, "-web.listen-address", address)
	go func() {
		if err := exporter.Run(); err != nil {
			errc <- fmt.Errorf("execution of node_exporter failed: %s", err)
		} else {
			errc <- nil
		}
	}()

	select {
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(100 * time.Millisecond):
	}

	go func(pid int, url string) {
		if err := queryExporter(url); err != nil {
			errc <- err
			return
		}
		proc, err := procfs.NewProc(pid)
		if err != nil {
			errc <- err
			return
		}
		fdsBefore, err := proc.FileDescriptors()
		if err != nil {
			errc <- err
			return
		}
		for i := 0; i < 5; i++ {
			if err := queryExporter(url); err != nil {
				errc <- err
				return
			}
		}
		fdsAfter, err := proc.FileDescriptors()
		if err != nil {
			errc <- err
			return
		}
		if want, have := len(fdsBefore), len(fdsAfter); want != have {
			errc <- fmt.Errorf("want %d open file descriptors after metrics scrape, have %d", want, have)
		}
		errc <- nil
	}(exporter.Process.Pid, fmt.Sprintf("http://%s/metrics", address))

	select {
	case err := <-errc:
		if exporter.Process != nil {
			exporter.Process.Kill()
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func queryExporter(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if want, have := resp.StatusCode, http.StatusOK; want != have {
		return fmt.Errorf("want /metrics status code %d, have %d", want, have)
	}
	return nil
}
