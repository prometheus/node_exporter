package model

type SystemInfo struct {
	Machine *Machine          `json:"machine,omitempty"`
	CPU     *CPU              `json:"cpu,omitempty"`
	Memory  *Memory           `json:"memory,omitempty"`
	Disk    *Disk             `json:"disk,omitempty"`
	GPU     *GPU              `json:"gpu,omitempty"`
	Net     *Net              `json:"net,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type Machine struct {
	Vendor       string `json:"vendor,omitempty"`
	Product      string `json:"product,omitempty"`
	Version      string `json:"version,omitempty"`
	Serial       string `json:"serial,omitempty"`
	UUID         string `json:"uuid,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	Kernel       string `json:"kernel,omitempty"`
	KernelArch   string `json:"kernel_arch,omitempty"`
	OS           string `json:"os,omitempty"`
	OSVersion    string `json:"os_version,omitempty"`
	Uptime       uint64 `json:"uptime,omitempty"`
	Type         string `json:"type,omitempty"`
	K8sNode      bool   `json:"k8s_node,omitempty"`
	BoardVendor  string `json:"board_vendor,omitempty"`
	BoardName    string `json:"board_name,omitempty"`
	BoardVersion string `json:"board_version,omitempty"`
	BoardSerial  string `json:"board_serial,omitempty"`
}

type CPU struct {
	Sockets int         `json:"sockets"`
	Cores   int         `json:"cores"`
	Threads int         `json:"threads"`
	Devices []CPUDevice `json:"devices,omitempty"`
}

type CPUDevice struct {
	ModelName string  `json:"model_name,omitempty"`
	VendorID  string  `json:"vendor_id,omitempty"`
	Cores     int     `json:"cores"`
	Mhz       float64 `json:"mhz,omitempty"`
	CacheKB   int     `json:"cache_kb,omitempty"`
}

type MemoryModule struct {
	Locator      string `json:"locator,omitempty"`
	BankLocator  string `json:"bank_locator,omitempty"`
	Size         string `json:"size,omitempty"`
	Type         string `json:"type,omitempty"`
	Speed        string `json:"speed,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Serial       string `json:"serial,omitempty"`
	PartNumber   string `json:"part_number,omitempty"`
}

type Memory struct {
	TotalBytes uint64         `json:"total_bytes"`
	Modules    []MemoryModule `json:"modules,omitempty"`
}

type DiskDevice struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Model      string `json:"model,omitempty"`
	Vendor     string `json:"vendor,omitempty"`
	Serial     string `json:"serial,omitempty"`
	SizeBytes  uint64 `json:"size_bytes,omitempty"`
	Mountpoint string `json:"mountpoint,omitempty"`
	FsType     string `json:"fs_type,omitempty"`
	UsedBytes  uint64 `json:"used_bytes,omitempty"`
}

type Disk struct {
	Devices []DiskDevice `json:"devices,omitempty"`
}

type GPUDevice struct {
	Index           int     `json:"index"`
	Vendor          string  `json:"vendor,omitempty"`
	Name            string  `json:"name,omitempty"`
	Serial          string  `json:"serial,omitempty"`
	UUID            string  `json:"uuid,omitempty"`
	Health          string  `json:"health,omitempty"`
	MemoryTotalMB   uint64  `json:"memory_total_mb,omitempty"`
	MemoryUsedMB    uint64  `json:"memory_used_mb,omitempty"`
	MemoryFreeMB    uint64  `json:"memory_free_mb,omitempty"`
	Utilization     float64 `json:"utilization,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
	PowerW          float64  `json:"power_w,omitempty"`
	DriverVersion   string  `json:"driver_version,omitempty"`
	FirmwareVersion string  `json:"firmware_version,omitempty"`
	// RuntimeMetrics reports whether memory/utilization/temperature/power
	// were actually read from the vendor tool. False means the device was
	// only enumerated via lspci (e.g. passthrough/vfio) and the runtime
	// fields are zero values — emitting them as 0 would be misleading.
	RuntimeMetrics bool `json:"runtime_metrics,omitempty"`
}

type GPU struct {
	Devices []GPUDevice `json:"devices,omitempty"`
}

type NetDevice struct {
	Name      string   `json:"name,omitempty"`
	Mac       string   `json:"mac,omitempty"`
	AddrsV4   []string `json:"addrs_v4,omitempty"`
	AddrsV6   []string `json:"addrs_v6,omitempty"`
	MTU       int      `json:"mtu,omitempty"`
	Up        bool     `json:"up"`
	Physical  bool     `json:"physical"`
	Master    string   `json:"master,omitempty"`
	Slaves    []string `json:"slaves,omitempty"`
	Vendor    string   `json:"vendor,omitempty"`
	Driver    string   `json:"driver,omitempty"`
	SpeedMbps int      `json:"speed_mbps,omitempty"`
}

type Net struct {
	Devices []NetDevice `json:"devices,omitempty"`
}
