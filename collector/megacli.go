// +build megacli

package collector

import (
	"bufio"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultMegaCli   = "megacli"
	adapterHeaderSep = "================"
)

var (
	driveTemperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "megacli_drive_temperature_celsius",
		Help:      "megacli: drive temperature",
	}, []string{"enclosure", "slot"})

	driveCounters = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "megacli_drive_count",
		Help:      "megacli: drive error and event counters",
	}, []string{"enclosure", "slot", "type"})

	drivePresence = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "megacli_adapter_disk_presence",
		Help:      "megacli: disk presence per adapter",
	}, []string{"type"})

	counters = []string{"Media Error Count", "Other Error Count", "Predictive Failure Count"}
)

func init() {
	Factories["megacli"] = NewMegaCliCollector
}

func parseMegaCliDisks(r io.ReadCloser) (map[int]map[int]map[string]string, error) {
	defer r.Close()
	stats := map[int]map[int]map[string]string{}
	scanner := bufio.NewScanner(r)

	curEnc := -1
	curSlot := -1
	for scanner.Scan() {
		var err error
		text := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(text, ":", 2)
		if len(parts) != 2 { // Adapter #X
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch {
		case key == "Enclosure Device ID":
			curEnc, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		case key == "Slot Number":
			curSlot, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		case curSlot != -1 && curEnc != -1:
			if _, ok := stats[curEnc]; !ok {
				stats[curEnc] = map[int]map[string]string{}
			}
			if _, ok := stats[curEnc][curSlot]; !ok {
				stats[curEnc][curSlot] = map[string]string{}
			}
			stats[curEnc][curSlot][key] = value
		}
	}
	return stats, nil
}

func parseMegaCliAdapter(r io.ReadCloser) (map[string]map[string]string, error) {
	defer r.Close()
	raidStats := map[string]map[string]string{}
	scanner := bufio.NewScanner(r)
	header := ""
	last := ""
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == adapterHeaderSep {
			header = last
			raidStats[header] = map[string]string{}
			continue
		}
		last = text
		if header == "" { // skip Adapter #X and separator
			continue
		}
		parts := strings.SplitN(text, ":", 2)
		if len(parts) != 2 { // these section never include anything we are interested in
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		raidStats[header][key] = value

	}
	return raidStats, nil
}

type megaCliCollector struct {
	config Config
	cli    string
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// RAID status through megacli.
func NewMegaCliCollector(config Config) (Collector, error) {
	cli := defaultMegaCli
	if config.Config["megacli_command"] != "" {
		cli = config.Config["megacli_command"]
	}

	c := megaCliCollector{
		config: config,
		cli:    cli,
	}

	if _, err := prometheus.RegisterOrGet(driveTemperature); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(driveCounters); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(drivePresence); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *megaCliCollector) Update() (updates int, err error) {
	au, err := c.updateAdapter()
	if err != nil {
		return au, err
	}
	du, err := c.updateDisks()
	return au + du, err
}

func (c *megaCliCollector) updateAdapter() (int, error) {
	cmd := exec.Command(c.cli, "-AdpAllInfo", "-aALL")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	stats, err := parseMegaCliAdapter(pipe)
	if err != nil {
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		return 0, err
	}

	updates := 0
	for k, v := range stats["Device Present"] {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return updates, err
		}
		drivePresence.WithLabelValues(k).Set(value)
		updates++
	}
	return updates, nil
}

func (c *megaCliCollector) updateDisks() (int, error) {
	cmd := exec.Command(c.cli, "-PDList", "-aALL")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	stats, err := parseMegaCliDisks(pipe)
	if err != nil {
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		return 0, err
	}

	updates := 0
	for enc, encStats := range stats {
		for slot, slotStats := range encStats {
			tStr := slotStats["Drive Temperature"]
			tStr = tStr[:strings.Index(tStr, "C")]
			t, err := strconv.ParseFloat(tStr, 64)
			if err != nil {
				return updates, err
			}

			encStr := strconv.Itoa(enc)
			slotStr := strconv.Itoa(slot)

			driveTemperature.WithLabelValues(encStr, slotStr).Set(t)
			updates++

			for _, c := range counters {
				counter, err := strconv.ParseFloat(slotStats[c], 64)
				if err != nil {
					return updates, err
				}

				driveCounters.WithLabelValues(encStr, slotStr, c).Set(counter)
				updates++
			}
		}
	}
	return updates, nil
}
