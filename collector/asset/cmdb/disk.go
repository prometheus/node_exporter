package cmdb

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
	"github.com/shirou/gopsutil/v3/disk"
)

type lsblkDevice struct {
	Name       string        `json:"name"`
	Model      string        `json:"model"`
	Vendor     string        `json:"vendor"`
	Serial     string        `json:"serial"`
	Size       flexUint      `json:"size"`
	Type       string        `json:"type"`
	Mountpoint string        `json:"mountpoint"`
	FSType     string        `json:"fstype"`
	Children   []lsblkDevice `json:"children,omitempty"`
}

type flexUint uint64

func (f *flexUint) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*f = flexUint(v)
	return nil
}

type lsblkOutput struct {
	BlockDevices []lsblkDevice `json:"blockdevices"`
}

func CollectDisk() (*model.Disk, error) {
	d := &model.Disk{Devices: []model.DiskDevice{}}

	collectLSBLK(d)
	collectPartitionsUsage(d)

	return d, nil
}

func collectLSBLK(d *model.Disk) {
	if !commandExists("lsblk") {
		return
	}

	out, err := runCmd("lsblk", "-b", "-J",
		"-o", "NAME,MODEL,VENDOR,SERIAL,SIZE,TYPE,MOUNTPOINT,FSTYPE")
	if err != nil {
		return
	}

	var parsed lsblkOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		return
	}

	for _, dev := range parsed.BlockDevices {
		walkLSBLK(dev, d)
	}
}

func walkLSBLK(dev lsblkDevice, d *model.Disk) {
	if uint64(dev.Size) == 0 {
		return
	}
	if dev.Type != "disk" {
		return
	}
	d.Devices = append(d.Devices, model.DiskDevice{
		Name:       dev.Name,
		Type:       dev.Type,
		Model:      dev.Model,
		Vendor:     dev.Vendor,
		Serial:     dev.Serial,
		SizeBytes:  uint64(dev.Size),
		Mountpoint: dev.Mountpoint,
		FsType:     dev.FSType,
	})
}

func collectPartitionsUsage(d *model.Disk) {
	parts, err := disk.Partitions(true)
	if err != nil {
		return
	}
	usageMap := map[string]uint64{}
	for _, p := range parts {
		if u, err := disk.Usage(p.Mountpoint); err == nil {
			usageMap[strings.TrimSpace(p.Device)] = u.Used
		}
	}
	for i, dev := range d.Devices {
		if dev.Mountpoint == "" {
			continue
		}
		devPath := "/dev/" + dev.Name
		if used, ok := usageMap[devPath]; ok {
			d.Devices[i].UsedBytes = used
		}
	}
}
