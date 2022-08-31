package percona_tests

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tklauser/go-sysconf"
)

const (
	repeatCount  = 5
	scrapesCount = 50
)

var doRun = flag.Bool("doRun", false, "")
var url = flag.String("url", "", "")

type StatsData struct {
	meanMs     float64
	stdDevMs   float64
	stdDevPerc float64

	meanHwm        float64
	stdDevHwmBytes float64
	stdDevHwmPerc  float64

	meanData        float64
	stdDevDataBytes float64
	stdDevDataPerc  float64
}

func TestPerformance(t *testing.T) {
	if !getBool(doRun) {
		t.Skip("For manual runs only through make")
		return
	}

	var updated, original *StatsData
	t.Run("upstream exporter", func(t *testing.T) {
		updated = doTestStats(t, repeatCount, scrapesCount, updatedExporterFileName)
	})

	t.Run("percona exporter", func(t *testing.T) {
		original = doTestStats(t, repeatCount, scrapesCount, oldExporterFileName)
	})

	printStats(original, updated)
}

func calculatePerc(base, updated float64) float64 {
	diff := base - updated
	diffPerc := float64(100) / base * diff
	diffPerc = diffPerc * -1

	return diffPerc
}

func doTestStats(t *testing.T, cnt int, size int, fileName string) *StatsData {
	var durations []float64
	var hwms []float64
	var datas []float64

	for i := 0; i < cnt; i++ {
		d, hwm, data, err := doTest(size, fileName)
		if !assert.NoError(t, err) {
			return nil
		}

		durations = append(durations, float64(d))
		hwms = append(hwms, float64(hwm))
		datas = append(datas, float64(data))
	}

	mean, _ := stats.Mean(durations)
	stdDev, _ := stats.StandardDeviation(durations)
	stdDev = float64(100) / mean * stdDev

	clockTicks, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		panic(err)
	}

	mean = mean * float64(1000) / float64(clockTicks) / float64(size)
	stdDevMs := stdDev / float64(100) * mean

	meanHwm, _ := stats.Mean(hwms)
	stdDevHwm, _ := stats.StandardDeviation(hwms)
	stdDevHwmPerc := float64(100) / meanHwm * stdDevHwm

	meanData, _ := stats.Mean(datas)
	stdDevData, _ := stats.StandardDeviation(datas)
	stdDevDataPerc := float64(100) / meanData * stdDevData

	st := StatsData{
		meanMs:     mean,
		stdDevMs:   stdDevMs,
		stdDevPerc: stdDev,

		meanHwm:        meanHwm,
		stdDevHwmBytes: stdDevHwm,
		stdDevHwmPerc:  stdDevHwmPerc,

		meanData:        meanData,
		stdDevDataBytes: stdDevData,
		stdDevDataPerc:  stdDevDataPerc,
	}

	//fmt.Printf("loop %dx%d: sample time: %.2fms [deviation ±%.2fms, %.1f%%]\n", cnt, scrapesCount, st.meanMs, st.stdDevMs, st.stdDevPerc)
	fmt.Printf("running %d scrapes %d times\n", size, cnt)
	fmt.Printf("CPU\t%.1fms [±%.1fms, %.1f%%]\n", st.meanMs, st.stdDevMs, st.stdDevPerc)
	fmt.Printf("HWM\t%.1fkB [±%.1f kB, %.1f%%]\n", st.meanHwm, st.stdDevHwmBytes, st.stdDevHwmPerc)
	fmt.Printf("Data\t%.1fkB [±%.1f kB, %.1f%%]\n", st.meanData, st.stdDevDataBytes, st.stdDevDataPerc)

	return &st
}

func doTest(iterations int, fileName string) (cpu, hwm, data int64, _ error) {
	cmd, port, collectOutput, err := launchExporter(fileName)
	if err != nil {
		return 0, 0, 0, err
	}

	total1 := getCPUTime(cmd.Process.Pid)

	for i := 0; i < iterations; i++ {
		_, err = tryGetMetrics(port)
		if err != nil {
			return 0, 0, 0, errors.Wrapf(err, "Failed to perform test iteration %d.%s", i, collectOutput())
		}

		time.Sleep(1 * time.Millisecond)
	}

	total2 := getCPUTime(cmd.Process.Pid)

	hwm, data = getCPUMem(cmd.Process.Pid)

	err = stopExporter(cmd, collectOutput)
	if err != nil {
		return 0, 0, 0, err
	}

	return total2 - total1, hwm, data, nil
}

func getCPUMem(pid int) (hwm, data int64) {
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(string(contents), "\n")

	for _, v := range lines {
		if strings.HasPrefix(v, "VmHWM") {
			val := strings.ReplaceAll(strings.ReplaceAll(strings.Split(v, ":\t")[1], " kB", ""), " ", "")
			hwm, _ = strconv.ParseInt(val, 10, 64)
			continue
		}
		if strings.HasPrefix(v, "VmData") {
			val := strings.ReplaceAll(strings.ReplaceAll(strings.Split(v, ":\t")[1], " kB", ""), " ", "")
			data, _ = strconv.ParseInt(val, 10, 64)
			continue
		}
	}

	return hwm, data
}

func getCPUTime(pid int) (total int64) {
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		numFields := len(fields)
		if numFields > 3 {
			i, err := strconv.ParseInt(fields[13], 10, 64)
			if err != nil {
				panic(err)
			}

			totalTime := i

			i, err = strconv.ParseInt(fields[14], 10, 64)
			if err != nil {
				panic(err)
			}

			totalTime += i

			total = totalTime

			return
		}
	}
	return
}

func printStats(original, updated *StatsData) {
	fmt.Println()
	fmt.Println("        \told\tnew\tdiff")
	fmt.Printf("CPU, ms \t%.1f\t%.1f\t%+.0f%%\n", original.meanMs, updated.meanMs, calculatePerc(original.meanMs, updated.meanMs))
	fmt.Printf("HWM, kB \t%.1f\t%.1f\t%+.0f%%\n", original.meanHwm, updated.meanHwm, calculatePerc(original.meanHwm, updated.meanHwm))
	fmt.Printf("DATA, kB\t%.1f\t%.1f\t%+.0f%%\n", original.meanData, updated.meanData, calculatePerc(original.meanData, updated.meanData))
	fmt.Println()
}
