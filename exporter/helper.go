package exporter

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func debug(name string, format string, a ...interface{}) {
	if *verbose {
		f := fmt.Sprintf("%s: %s", name, format)
		log.Printf(f, a...)
	}
}

func splitToInts(str string, sep string) (ints []int, err error) {
	for _, part := range strings.Split(str, sep) {
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Could not split '%s' because %s is no int: %s", str, part, err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}
