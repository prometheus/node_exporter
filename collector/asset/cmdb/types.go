package cmdb

// Re-export the data model types at the cmdb package level so consumers (the
// asset_* Prometheus collectors) can spell the collected value types as
// cmdb.Machine, cmdb.CPU, ... without importing the model sub-package. These
// are type aliases, so *cmdb.Machine is identical to *model.Machine: the values
// returned by CollectMachine/CollectCPU/etc. can be used interchangeably.
import "github.com/prometheus/node_exporter/collector/asset/cmdb/model"

type (
	Machine = model.Machine
	CPU     = model.CPU
	Memory  = model.Memory
	Disk    = model.Disk
	GPU     = model.GPU
	Net     = model.Net
)
