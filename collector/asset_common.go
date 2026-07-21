// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

// Asset collector namespace. Every metric exposed by the asset_* collectors is
// built with this prefix, yielding siliconflow_asset_<name>.
const assetNamespace = "siliconflow_asset"

// assetCacheTTL bounds how often each asset_* collector re-runs its cmdb
// collection. Asset inventory is relatively static (machine/cpu/memory/disk/net
// hardware doesn't change at runtime; gpu runtime metrics drift slowly), so
// collectors cache the cmdb.Collect* result and serve it on every scrape until
// the TTL elapses, then refresh. This keeps scrape latency low (no dmidecode /
// nvidia-smi / lsblk shell-out per scrape) and matches a "one snapshot per day"
// inventory cadence. Set to 0 to disable caching and collect on every scrape.
var assetCacheTTL = kingpin.Flag(
	"collector.asset.cache-ttl",
	"Time-to-live for cached asset collection results. Only applies to "+
		"purely-static collectors (asset_cpu, asset_memory); collectors with "+
		"realtime fields (asset_machine uptime, asset_disk used_bytes, "+
		"asset_gpu utilization/temperature/power/memory, asset_net speed) "+
		"always collect on every scrape. Set to 0 to disable caching entirely.",
).Default("24h").Duration()

// assetUUIDFilePath is the persistent UUID written by --generate-uuid and read
// by every asset_* collector to label its metrics with asset_uuid. Hardcoded by
// design; change here to relocate. Exported as AssetUUIDFilePath so the main
// binary can reference it in the --generate-uuid flag help text.
const assetUUIDFilePath = "/var/lib/siliconflow_asset/uuid"

// AssetUUIDFilePath is the filesystem path of the persistent asset UUID file.
const AssetUUIDFilePath = assetUUIDFilePath

// assetUUIDLabel is the label name carrying the persistent asset UUID on every
// siliconflow_asset_* metric. Named asset_uuid (not "uuid") to avoid confusing
// it with the SMBIOS product UUID exposed as the machine_uuid label of
// siliconflow_asset_machine_info.
const assetUUIDLabel = "asset_uuid"

// readAssetUUID reads the persistent asset UUID written by --generate-uuid.
// It is called by each asset_* collector at scrape time (not at construction
// time) so the UUID may be generated after the exporter has started. On failure
// the collector returns the error: node_scrape_collector_success{collector=...}
// is set to 0 and an error is logged, so a missing UUID never produces metrics
// with an empty asset_uuid label.
func readAssetUUID() (string, error) {
	b, err := os.ReadFile(assetUUIDFilePath)
	if err != nil {
		return "", fmt.Errorf("read asset uuid %s: %w (generate it first with --generate-uuid)", assetUUIDFilePath, err)
	}
	return strings.TrimSpace(string(b)), nil
}

// assetLabel sanitizes a string for use as a Prometheus label value. DMI/SMBIOS
// and sysfs strings occasionally contain invalid UTF-8 sequences; replace them
// so the text-format exposition stays valid (mirrors the dmi collector).
func assetLabel(s string) string {
	return strings.ToValidUTF8(s, "?")
}

// assetBool renders a bool as a stable "true"/"false" label value.
func assetBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// assetCache is a thread-safe TTL cache for a single collected value. Each
// asset_* collector holds one instance for its cmdb result. On a cache miss the
// given fetch function runs and its result is cached for ttl; while the cache is
// fresh, scrapes serve the cached value without re-running dmidecode /
// nvidia-smi / lsblk / lspci. If a refresh fails and a previous value exists,
// the stale value is served so a transient collection failure doesn't drop
// inventory metrics; only when no value has ever been collected is the error
// propagated (setting node_scrape_collector_success=0).
type assetCache[T any] struct {
	mu       sync.Mutex
	value    T
	fetched  time.Time
	hasValue bool
}

func (c *assetCache[T]) get(ttl time.Duration, fetch func() (T, error)) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl > 0 && c.hasValue && time.Since(c.fetched) < ttl {
		return c.value, nil
	}
	v, err := fetch()
	if err != nil {
		if c.hasValue {
			// Serve stale cache on transient failure rather than dropping
			// inventory metrics. The cache timestamp is NOT advanced, so the
			// next scrape will retry the refresh.
			return c.value, nil
		}
		var zero T
		return zero, err
	}
	c.value = v
	c.fetched = time.Now()
	c.hasValue = true
	return v, nil
}
