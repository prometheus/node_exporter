package dbus

import (
	"errors"
	"os/exec"
)

func getSessionBusPlatformAddress() (string, error) {
	cmd := exec.Command("launchctl", "getenv", "DBUS_LAUNCHD_SESSION_BUS_SOCKET")
	b, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	if len(b) == 0 {
		return "", errors.New("dbus: couldn't determine address of session bus")
	}

	return "unix:path=" + string(b[:len(b)-1]), nil
}
