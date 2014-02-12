package main

import (
	"flag"
	"log"

	"github.com/prometheus/node_exporter/exporter"
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
	log.Printf("Registered collectors:")
	for _, c := range exporter.Collectors {
		log.Print(" - ", c.Name())
	}
	exporter.Loop()
}
