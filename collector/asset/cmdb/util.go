package cmdb

import (
	"os"
	"os/exec"
	"strings"
)

const dmiPath = "/sys/class/dmi/id/"

func readDMIFile(name string) string {
	b, err := os.ReadFile(dmiPath + name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func runCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func readSysFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}
