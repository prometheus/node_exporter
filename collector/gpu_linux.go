// +build !nogpu

package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/common/log"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
}

type gpuInfo struct {
	TotalMem    float64
	UsedMem     float64
	FreeMem     float64
	Utilization float64
	Temp        float64
	Host        string
	UUID        string
	Count       string
	Types       string
}

type gpuCache struct{}

func (this gpuCache) Stat() ([]gpuInfo, error) {
	var (
		result []gpuInfo
		// err    error
	)
	command := "nvidia-smi -q|grep -E 'Minor Number|UUID|GPU Current Temp|Gpu|Total|Used|Free|Product Name'|grep -v 'Used GPU Memory'|cut -d ':' -f2|sed 's/^[[:space:]]//g'"
	data, err := execCommand(command)
	if err != nil {
		return nil, err
	}

	result = []gpuInfo{}
	var tmp gpuInfo
	for n, x := range strings.Split(data, "\n") {
		// fmt.Println(n, n%15, n/15, x)
		if n%15 == 0 && n != 0 {
			result = append(result, tmp)
			fmt.Println("n = 0", x)
			tmp = gpuInfo{}
			tmp.Host = hostname
			tmp.Types = strings.TrimSpace(x)
		} else if n%15 == 0 && n == 0 {
			fmt.Println("n = 0", x)
			tmp = gpuInfo{}
			tmp.Host = hostname
			tmp.Types = strings.TrimSpace(x)
		} else if n%15 == 1 {
			tmp.UUID = strings.TrimSpace(x)
		} else if n%15 == 2 {
			tmp.Count = strings.TrimSpace(x)
		} else if n%15 == 3 {
			tmp.TotalMem = strings.TrimSpace(strings.Split(x, " ")[0])
		} else if n%15 == 4 {
			tmp.UsedMem = strings.TrimSpace(strings.Split(x, " ")[0])
		} else if n%15 == 5 {
			tmp.FreeMem = strings.TrimSpace(strings.Split(x, " ")[0])
		} else if n%15 == 9 {
			tmp.Utilization = strings.TrimSpace(strings.Split(x, " ")[0])
		} else if n%15 == 14 {
			fmt.Println("temp", x)
			tmp.Temp = strings.TrimSpace(strings.Split(x, " ")[0])
		}
	}
	return result, nil
}

func execCommand(cmd string) (string, error) {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Stdout = &out
	pipeline.Stderr = &stderr
	err := pipeline.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func test() {
	tmp := gpuCache{}
	x, err := tmp.Stat()
	if err != nil {
		panic(err)
	}

	xxx, _ := json.Marshal(x)
	fmt.Println(string(xxx))
}
