package cmdb

import (
	"os"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
	"github.com/shirou/gopsutil/v3/host"
)

var vmVendors = map[string]string{
	"vmware":                "VMware",
	"vmware, inc.":          "VMware",
	"qemu":                  "QEMU",
	"kvm":                   "KVM",
	"xen":                   "Xen",
	"microsoft corporation": "Hyper-V",
	"innotek gmbh":          "VirtualBox",
	"parallels":             "Parallels",
	"openvz":                "OpenVZ",
	"bochs":                 "Bochs",
	"oracle corporation":    "VirtualBox",
	"amazon ec2":            "Xen",
}

func CollectMachine() (*model.Machine, error) {
	m := &model.Machine{}

	if prod, err := ghw.Product(); err == nil && prod != nil {
		m.Vendor = prod.Vendor
		m.Product = prod.Name
		m.Version = prod.Version
		m.Serial = firstNonEmpty(prod.SerialNumber, readDMIFile("product_serial"))
		m.UUID = prod.UUID
	} else {
		m.Vendor = readDMIFile("sys_vendor")
		m.Product = readDMIFile("product_name")
		m.Version = readDMIFile("product_version")
		m.Serial = readDMIFile("product_serial")
		m.UUID = readDMIFile("product_uuid")
	}

	var virtSystem, virtRole string
	if info, err := host.Info(); err == nil {
		m.Hostname = info.Hostname
		m.OS = firstNonEmpty(info.Platform, info.OS)
		m.OSVersion = info.PlatformVersion
		m.Kernel = info.KernelVersion
		m.KernelArch = info.KernelArch
		m.Uptime = info.Uptime
		virtSystem = info.VirtualizationSystem
		virtRole = info.VirtualizationRole
	}

	m.Type = detectMachineType(m.Vendor, m.Product, virtSystem, virtRole)
	m.K8sNode = detectK8sNode()

	if bb, err := ghw.Baseboard(); err == nil && bb != nil {
		m.BoardVendor = bb.Vendor
		m.BoardName = bb.Product
		m.BoardVersion = bb.Version
		m.BoardSerial = firstNonEmpty(bb.SerialNumber, readDMIFile("board_serial"))
	} else {
		m.BoardVendor = readDMIFile("board_vendor")
		m.BoardName = readDMIFile("board_name")
		m.BoardVersion = readDMIFile("board_version")
		m.BoardSerial = readDMIFile("board_serial")
	}

	return m, nil
}

// detectK8sNode 通过 kubelet 进程或 /var/lib/kubelet 判定本机是否为 K8s 节点。
// 容器网络网卡会被 net 采集器排除,这里仅做节点类型标记。
func detectK8sNode() bool {
	if commandExists("kubelet") {
		if out, err := runCmd("pgrep", "-x", "kubelet"); err == nil && strings.TrimSpace(out) != "" {
			return true
		}
	}
	if info, err := os.Stat("/var/lib/kubelet"); err == nil && info.IsDir() {
		return true
	}
	return false
}

func detectMachineType(vendor, product, virtSystem, virtRole string) string {
	if virtRole == "guest" && virtSystem != "" {
		return "virtual"
	}

	v := strings.ToLower(strings.TrimSpace(vendor))
	p := strings.ToLower(strings.TrimSpace(product))
	for sig := range vmVendors {
		if strings.Contains(v, sig) || strings.Contains(p, sig) {
			return "virtual"
		}
	}

	if virtRole == "guest" {
		return "virtual"
	}

	return "physical"
}
