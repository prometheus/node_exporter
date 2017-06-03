// +build !darwin

package iostat

import (
	"errors"
)

// ReadDriveStats returns statictics of each of the drives.
func ReadDriveStats() ([]*DriveStats, error) {
	return nil, errors.New("not implement")
}
