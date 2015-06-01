// +build !nonetdev

package collector

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <net/if.h>
*/
import "C"

const (
	netDevSubsystem = "network"
)

type netDevCollector struct {
	config  Config
	metrics map[string]*prometheus.CounterVec
}

func init() {
	Factories["netdev"] = NewNetDevCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device stats.
func NewNetDevCollector(config Config) (Collector, error) {
	return &netDevCollector{
		config:  config,
		metrics: map[string]*prometheus.CounterVec{},
	}, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netDev, err := getNetDevStats()
	if err != nil {
		return fmt.Errorf("Couldn't get netstats: %s", err)
	}
	for direction, devStats := range netDev {
		for dev, stats := range devStats {
			for t, value := range stats {
				key := direction + "_" + t
				if _, ok := c.metrics[key]; !ok {
					c.metrics[key] = prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Namespace: Namespace,
							Subsystem: netDevSubsystem,
							Name:      key,
							Help:      fmt.Sprintf("%s %s from getifaddrs().", t, direction),
						},
						[]string{"device"},
					)
				}
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("Invalid value %s in netstats: %s", value, err)
				}
				c.metrics[key].WithLabelValues(dev).Set(v)
			}
		}
	}
	for _, m := range c.metrics {
		m.Collect(ch)
	}
	return err
}

func getNetDevStats() (map[string]map[string]map[string]string, error) {
	netDev := map[string]map[string]map[string]string{}
	netDev["transmit"] = map[string]map[string]string{}
	netDev["receive"] = map[string]map[string]string{}

	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed!")
	}
	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family == C.AF_LINK {
			receive := map[string]string{}
			transmit := map[string]string{}
			//var data *C.struct_if_data
			data := (*C.struct_if_data)(ifa.ifa_data)
			receive["packets"] = strconv.Itoa(int(data.ifi_ipackets))
			transmit["packets"] = strconv.Itoa(int(data.ifi_opackets))
			receive["errs"] = strconv.Itoa(int(data.ifi_ierrors))
			transmit["errs"] = strconv.Itoa(int(data.ifi_oerrors))
			receive["bytes"] = strconv.Itoa(int(data.ifi_ibytes))
			transmit["bytes"] = strconv.Itoa(int(data.ifi_obytes))
			receive["multicast"] = strconv.Itoa(int(data.ifi_imcasts))
			transmit["multicast"] = strconv.Itoa(int(data.ifi_omcasts))
			receive["drop"] = strconv.Itoa(int(data.ifi_iqdrops))
			//transmit["drops"] = strconv.Itoa(int(data.ifi_oqdrops))

			netDev["receive"][C.GoString(ifa.ifa_name)] = receive
			netDev["transmit"][C.GoString(ifa.ifa_name)] = transmit
		}
	}
	defer C.freeifaddrs(ifap)

	return netDev, nil
}
