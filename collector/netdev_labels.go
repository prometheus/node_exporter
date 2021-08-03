package collector

import (
	"io/ioutil"
	"runtime"
	"strings"
)

type label struct {
	key   string
	value string
}

type labels []label

func (l labels) keys() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.key)
	}

	return ret
}

func (l labels) values() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.value)
	}

	return ret
}

func getLabelsFromIfAlias(ifName string) labels {
	if runtime.GOOS != "linux" {
		return nil
	}

	if !labelsFromIfAlias {
		return nil
	}

	ifAliasBytes, err := ioutil.ReadFile("/sys/class/net/" + ifName + "/ifalias")
	if err != nil {
		return nil
	}

	ifAlias := strings.TrimSpace(string(ifAliasBytes))
	keyValueStrings := strings.Split(ifAlias, ",")
	ret := make(labels, 0, len(keyValueStrings))

	for _, kv := range keyValueStrings {
		parts := strings.Split(kv, "=")
		if len(parts) != 2 {
			continue
		}

		ret = append(ret, label{
			key:   parts[0],
			value: parts[1],
		})
	}

	return ret
}
