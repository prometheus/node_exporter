// +build !nogpu

package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
}

type gpuCache struct{}

/*
nvidia-smi -q|grep -E 'Minor Number|UUID|GPU Current Temp|Gpu|Total|Used|Free'|cut -d ':' -f2|awk '{print $1}'
GPU-1111111111111111111111
0
11178
0
11178
256
2
254
0
N/A
N/A
N/A
N/A
31
GPU-a2222222222222222222222222222
1
11178
0
11178
256
2
254
0
N/A
N/A
N/A
N/A
28
*/
func (this gpuCache) Stat() ([]gpuInfo, error) {
	var (
		result []gpuInfo
		// err    error
	)
	command := "nvidia-smi -q|grep -E 'Minor Number|UUID|GPU Current Temp|Gpu|Total|Used|Free'|grep -v 'Used GPU Memory'|cut -d ':' -f2|awk '{print $1}'"
	data, err := execCommand(command)
	if err != nil {
		return nil, err
	}

	result = []gpuInfo{}
	var tmp gpuInfo
	for n, x := range strings.Split(data, "\n") {
		// fmt.Println(n, n%14, n/14, x)
		log.Debug(n, n%14, n/14, x)
		if n%14 == 0 && n != 0 {
			result = append(result, tmp)
			tmp = gpuInfo{}
			tmp.Host = hostname
			tmp.UUID = strings.TrimSpace(x)
		} else if n%14 == 0 && n == 0 {
			tmp = gpuInfo{}
			tmp.Host = hostname
			tmp.UUID = strings.TrimSpace(x)
		} else if n%14 == 1 {
			tmp.Count = strings.TrimSpace(x)
		} else if n%14 == 2 {
			tmp.TotalMem, _ = strconv.ParseFloat(strings.TrimSpace(x), 64)
		} else if n%14 == 3 {
			tmp.UsedMem, _ = strconv.ParseFloat(strings.TrimSpace(x), 64)
		} else if n%14 == 4 {
			tmp.FreeMem, _ = strconv.ParseFloat(strings.TrimSpace(x), 64)
		} else if n%14 == 8 {
			tmp.Utilization, _ = strconv.ParseFloat(strings.TrimSpace(x), 64)
		} else if n%14 == 13 {
			tmp.Temp, _ = strconv.ParseFloat(strings.TrimSpace(x), 64)
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
