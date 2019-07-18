// +build !nogpu
// +build linux freebsd openbsd darwin,amd64 dragonfly

package collector

const (
	gpuCollectorSubsystem = "gpu"
)

var (
	gpuLabelNames = []string{"hostname", "id", "uuid"}
)
