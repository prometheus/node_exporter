// +build !nopending_updates

package collector

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultCheckCommand = "/usr/lib/update-notifier/apt-check"
)

var (
	checkCommand   = flag.String("checkUpdatesCommand", defaultCheckCommand, "command to run, returning pending updates as <all>;<security>")
	pendingUpdates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "pending_updates",
			Help:      "Number of pending updates.",
		}, []string{"type"})
)

type pendingUpdatesCollector struct{}

func init() {
	Factories["pending_updates"] = NewPendingUpdatesCollector
}

// NewPendingUpdatesCollector returns a newly allocated pendingUpdatesCollector.
// It exposes the number of pending package updates.
func NewPendingUpdatesCollector(config Config) (Collector, error) {
	c := pendingUpdatesCollector{}
	if _, err := prometheus.RegisterOrGet(pendingUpdates); err != nil {
		return nil, err
	}
	return &c, nil
}

// Update gathers pending updates, implements Collector interface
func (c *pendingUpdatesCollector) Update() (int, error) {
	out, err := exec.Command(*checkCommand).CombinedOutput()
	if err != nil {
		return 0, err
	}
	fields := strings.Split(string(out), ";")
	if len(fields) < 2 {
		return 0, fmt.Errorf("Expected %s to return pending updates as <all>;<security> but got: %s", *checkCommand, string(out))
	}
	allUpdates, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, err
	}
	secUpdates, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, err
	}
	pendingUpdates.WithLabelValues("all").Set(float64(allUpdates))
	pendingUpdates.WithLabelValues("security").Set(float64(secUpdates))
	return 2, nil
}
