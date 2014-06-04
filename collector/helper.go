package collector

import (
	"fmt"
	"strconv"
	"strings"
)

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
