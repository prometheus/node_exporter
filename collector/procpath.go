package collector

import (
	"flag"
	"path"

	"github.com/prometheus/procfs"
)

var (
	// The path of the proc filesystem.
	procPath = flag.String("collector.procfs", procfs.DefaultMountPoint, "procfs mountpoint.")
)

func procFilePath(name string) string {
	return path.Join(*procPath, name)
}
