package cmdb

import (
	"sort"
	"strconv"

	"github.com/prometheus/node_exporter/collector/asset/cmdb/model"
	"github.com/shirou/gopsutil/v3/cpu"
)

func CollectCPU() (*model.CPU, error) {
	c := &model.CPU{}

	infos, _ := cpu.Info()
	logical, _ := cpu.Counts(true)
	physical, _ := cpu.Counts(false)

	if logical == 0 {
		logical = len(infos)
	}
	c.Threads = logical

	type socket struct {
		coreIDs map[string]struct{}
		threads int
		hasCore bool
		example cpu.InfoStat
	}
	socks := map[string]*socket{}
	var order []string
	for _, ci := range infos {
		pid := ci.PhysicalID
		if pid == "" {
			pid = "0"
		}
		s, ok := socks[pid]
		if !ok {
			s = &socket{coreIDs: map[string]struct{}{}, example: ci}
			socks[pid] = s
			order = append(order, pid)
		}
		s.threads++
		if ci.CoreID != "" {
			s.hasCore = true
			s.coreIDs[ci.CoreID] = struct{}{}
		}
	}

	// 退化场景:每个逻辑线程都被报告为独立插槽(部分虚拟机给每个 vCPU
	// 分配独立 physical id)。此时拓扑无意义,折叠为单路。
	if len(socks) > 1 && len(socks) == len(infos) {
		var ex cpu.InfoStat
		if len(infos) > 0 {
			ex = infos[0]
		}
		socks = map[string]*socket{"0": {coreIDs: map[string]struct{}{}, example: ex}}
		order = []string{"0"}
	}

	sort.Slice(order, func(i, j int) bool {
		a, _ := strconv.Atoi(order[i])
		b, _ := strconv.Atoi(order[j])
		return a < b
	})

	totalCores := 0
	c.Devices = make([]model.CPUDevice, 0, len(order))
	for _, pid := range order {
		s := socks[pid]
		cores := 0
		if s.hasCore {
			cores = len(s.coreIDs)
		}
		totalCores += cores
		c.Devices = append(c.Devices, model.CPUDevice{
			ModelName: s.example.ModelName,
			VendorID:  s.example.VendorID,
			Cores:     cores,
			Mhz:       s.example.Mhz,
			CacheKB:   int(s.example.CacheSize),
		})
	}

	if totalCores > 0 {
		c.Sockets = len(order)
		c.Cores = totalCores
	} else {
		// 无法从 cpuinfo 获取 core_id 拓扑(虚拟机/容器常见):
		// 报告单路,核数取物理核数,缺失时退化为逻辑核数。
		cores := physical
		if cores == 0 {
			cores = logical
		}
		c.Sockets = 1
		c.Cores = cores
		dev := model.CPUDevice{Cores: cores}
		if len(infos) > 0 {
			dev.ModelName = infos[0].ModelName
			dev.VendorID = infos[0].VendorID
			dev.Mhz = infos[0].Mhz
			dev.CacheKB = int(infos[0].CacheSize)
		}
		c.Devices = append(c.Devices[:0], dev)
	}

	return c, nil
}
