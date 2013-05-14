package main

import (
	"flag"
	"github.com/prometheus/node_exporter/exporter"
	"log"
)

var (
	configFile = flag.String("config", "node_exporter.conf", "config file.")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {
	flag.Parse()

	exporter, err := exporter.New(*configFile)
	if err != nil {
		log.Fatalf("Couldn't instantiate exporter: %s", err)
	}
	exporter.Loop()
}
