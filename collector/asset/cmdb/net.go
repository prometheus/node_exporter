package cmdb

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
	gpnet "github.com/shirou/gopsutil/v3/net"
)

// virtualPrefixes 为容器/K8s 网络与回环等非资产网卡前缀,采集时整体排除。
// 采用方案 A:K8s 容器网络完全不进 CMDB,仅保留物理网卡与宿主机网络配置。
var virtualPrefixes = []string{
	"lo",        // loopback
	"docker",    // docker0 默认桥
	"br-",       // docker 自定义网络桥
	"veth",      // 容器 veth pair
	"cni",       // CNI 桥
	"flannel",   // flannel overlay
	"calico",    // calico
	"cilium",    // cilium
	"tunl",      // calico tunnel
	"genev",     // geneve tunnel
	"kube-ipvs", // kube-proxy ipvs dummy
	"ovn",       // ovn-kubernetes
	"nodelocaldns",
}

func CollectNet() (*model.Net, error) {
	n := &model.Net{Devices: []model.NetDevice{}}

	ifaces, err := gpnet.Interfaces()
	if err != nil {
		return n, err
	}

	// 先扫一遍,建立 slave -> bond 映射和 bond -> slaves 映射,
	// 用于后续在物理网卡上标记 master、在 bond 上列出 slaves。
	bondOf := map[string]string{}     // slaveName -> bondName
	slavesOf := map[string][]string{} // bondName -> []slaveName
	for _, iface := range ifaces {
		if !isBondNIC(iface.Name) {
			continue
		}
		slaves := readBondSlaves(iface.Name)
		slavesOf[iface.Name] = slaves
		for _, s := range slaves {
			bondOf[strings.TrimSpace(s)] = iface.Name
		}
	}

	for _, iface := range ifaces {
		if isVirtualInterface(iface.Name) {
			continue
		}

		physical := isPhysicalNIC(iface.Name)
		bond := isBondNIC(iface.Name)

		// 仅保留: 物理以太网卡(mlx5_core 等驱动的 IPoIB 接口虽挂在
		// /sys/class/net/<name>/device 上但属于 InfiniBand 虚拟 L3 口,
		// 不计入以太网资产) 与 bond 聚合口。bridge/vlan/tap/ovs 等全部丢弃。
		if physical && isInfiniBandNIC(iface.Name) {
			continue
		}
		if !physical && !bond {
			continue
		}

		dev := model.NetDevice{
			Name:     iface.Name,
			Mac:      strings.ToLower(strings.TrimSpace(iface.HardwareAddr)),
			MTU:      iface.MTU,
			Up:       isUp(iface.Flags),
			Physical: physical,
		}
		dev.AddrsV4, dev.AddrsV6 = classifyAddrs(iface.Addrs)

		if physical {
			dev.Vendor = readSysNetVendor(iface.Name)
			dev.Driver = readSysNetDriver(iface.Name)
			dev.SpeedMbps = readSysNetSpeed(iface.Name)
			// 标记从属的 bond(ens12f0 -> bond0);未加入 bond 的为空,单独显示。
			dev.Master = bondOf[iface.Name]
		} else if bond {
			dev.Driver = "bonding"
			dev.Slaves = slavesOf[iface.Name]
		}

		n.Devices = append(n.Devices, dev)
	}

	return n, nil
}

// readBondSlaves 读取 /sys/class/net/<bond>/bonding/slaves,返回该 bond
// 下挂的物理从口列表(空格分隔)。
func readBondSlaves(name string) []string {
	s := strings.TrimSpace(readSysFile(filepath.Join("/sys/class/net", name, "bonding", "slaves")))
	if s == "" {
		return nil
	}
	return strings.Fields(s)
}

func isVirtualInterface(name string) bool {
	low := strings.ToLower(name)
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(low, p) {
			return true
		}
	}
	return false
}

// isPhysicalNIC 通过 /sys/class/net/<name>/device 是否存在判定是否有
// PCI/USB 等总线背板(物理网卡)。bonds/vlans/bridges 没有该路径。
func isPhysicalNIC(name string) bool {
	_, err := os.Stat(filepath.Join("/sys/class/net", name, "device"))
	return err == nil
}

// isBondNIC 判定是否为 bond 聚合口。bond master 在
// /sys/class/net/<name>/bonding 下有 mode/slaves 等属性。
func isBondNIC(name string) bool {
	fi, err := os.Stat(filepath.Join("/sys/class/net", name, "bonding"))
	return err == nil && fi.IsDir()
}

// isInfiniBandNIC 通过 /sys/class/net/<name>/type 判定接口链路层类型。
// type=32 (ARPHRD_INFINIBAND) 表示 IPoIB,是跑在 IB 硬件之上的虚拟 L3 口,
// 不计入以太网资产。type=1 (ARPHRD_ETHER) 才是以太网。
func isInfiniBandNIC(name string) bool {
	t := strings.TrimSpace(readSysFile(filepath.Join("/sys/class/net", name, "type")))
	return t == "32"
}

func isUp(flags []string) bool {
	for _, f := range flags {
		if strings.ToLower(f) == "up" {
			return true
		}
	}
	return false
}

// classifyAddrs 将地址分为 v4/v6,排除 link-local(169.254/16 与 fe80::/10)
// 以及带 zone 的 IPv6 地址(%iface 后缀)。
func classifyAddrs(addrs []gpnet.InterfaceAddr) (v4, v6 []string) {
	for _, a := range addrs {
		ipStr := a.Addr
		if i := strings.Index(ipStr, "%"); i >= 0 {
			ipStr = ipStr[:i]
		}
		if slash := strings.Index(ipStr, "/"); slash >= 0 {
			ipStr = ipStr[:slash]
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}
		if ip.IsLinkLocalUnicast() {
			continue
		}
		if ip.To4() != nil {
			v4 = append(v4, ip.String())
		} else {
			v6 = append(v6, ip.String())
		}
	}
	return
}

func readSysNetVendor(name string) string {
	return strings.TrimSpace(readSysFile(filepath.Join("/sys/class/net", name, "device", "vendor")))
}

func readSysNetDriver(name string) string {
	link, err := os.Readlink(filepath.Join("/sys/class/net", name, "device", "driver"))
	if err != nil {
		return ""
	}
	return filepath.Base(link)
}

func readSysNetSpeed(name string) int {
	s := strings.TrimSpace(readSysFile(filepath.Join("/sys/class/net", name, "speed")))
	if s == "" {
		return 0
	}
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}
