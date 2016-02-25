package collector

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// zpool metrics

func (c *zfsCollector) parseZpoolOutput(reader io.Reader, handler func(string, string, float64)) (err error) {

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {

		fields := strings.Fields(scanner.Text())
		if len(fields) != 4 {
			return fmt.Errorf("Unexpected output of zpool command")
		}

		valueString := fields[2]
		switch {
		case strings.HasSuffix(fields[2], "%"):
			percentage := strings.TrimSuffix(fields[2], "%")
			valueString = "0." + percentage
		case strings.HasSuffix(fields[2], "x"):
			valueString = strings.TrimSuffix(fields[2], "x")
		}

		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return err
		}
		handler(fields[0], fields[1], value)

	}
	return scanner.Err()

}
