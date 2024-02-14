// Copyright 2016 The Prometheus Authors
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

//go:build !nozfs
// +build !nozfs

package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArcstatsParsing(t *testing.T) {
	arcstatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/arcstats")
	if err != nil {
		t.Fatal(err)
	}
	defer arcstatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(arcstatsFile, "arcstats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.arcstats.hits") {
			return
		}

		handlerCalled = true

		if v != uint64(8772612) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Arcstats parsing handler was not called for some expected sysctls")
	}
}

func TestZfetchstatsParsing(t *testing.T) {
	zfetchstatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/zfetchstats")
	if err != nil {
		t.Fatal(err)
	}
	defer zfetchstatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(zfetchstatsFile, "zfetchstats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.zfetchstats.hits") {
			return
		}

		handlerCalled = true

		if v != uint64(7067992) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Zfetchstats parsing handler was not called for some expected sysctls")
	}
}

func TestZilParsing(t *testing.T) {
	zilFile, err := os.Open("fixtures/proc/spl/kstat/zfs/zil")
	if err != nil {
		t.Fatal(err)
	}
	defer zilFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(zilFile, "zil", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.zil.zil_commit_count") {
			return
		}

		handlerCalled = true

		if v != uint64(10) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Zil parsing handler was not called for some expected sysctls")
	}
}

func TestVdevCacheStatsParsing(t *testing.T) {
	vdevCacheStatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/vdev_cache_stats")
	if err != nil {
		t.Fatal(err)
	}
	defer vdevCacheStatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(vdevCacheStatsFile, "vdev_cache_stats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.vdev_cache_stats.delegations") {
			return
		}

		handlerCalled = true

		if v != uint64(40) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("VdevCacheStats parsing handler was not called for some expected sysctls")
	}
}

func TestXuioStatsParsing(t *testing.T) {
	xuioStatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/xuio_stats")
	if err != nil {
		t.Fatal(err)
	}
	defer xuioStatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(xuioStatsFile, "xuio_stats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.xuio_stats.onloan_read_buf") {
			return
		}

		handlerCalled = true

		if v != uint64(32) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("XuioStats parsing handler was not called for some expected sysctls")
	}
}

func TestFmParsing(t *testing.T) {
	fmFile, err := os.Open("fixtures/proc/spl/kstat/zfs/fm")
	if err != nil {
		t.Fatal(err)
	}
	defer fmFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(fmFile, "fm", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.fm.erpt-dropped") {
			return
		}

		handlerCalled = true

		if v != uint64(18) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Fm parsing handler was not called for some expected sysctls")
	}
}

func TestDmuTxParsing(t *testing.T) {
	dmuTxFile, err := os.Open("fixtures/proc/spl/kstat/zfs/dmu_tx")
	if err != nil {
		t.Fatal(err)
	}
	defer dmuTxFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(dmuTxFile, "dmu_tx", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.dmu_tx.dmu_tx_assigned") {
			return
		}

		handlerCalled = true

		if v != uint64(3532844) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("DmuTx parsing handler was not called for some expected sysctls")
	}
}

func TestZpoolParsing(t *testing.T) {
	zpoolPaths, err := filepath.Glob("fixtures/proc/spl/kstat/zfs/*/io")
	if err != nil {
		t.Fatal(err)
	}

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	for _, zpoolPath := range zpoolPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			t.Fatal(err)
		}

		err = c.parsePoolProcfsFile(file, zpoolPath, func(poolName string, s zfsSysctl, v uint64) {
			if s != zfsSysctl("kstat.zfs.misc.io.nread") {
				return
			}

			handlerCalled = true

			if v != uint64(1884160) && v != uint64(2826240) {
				t.Fatalf("Incorrect value parsed from procfs data %v", v)
			}

		})
		file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	if !handlerCalled {
		t.Fatal("Zpool parsing handler was not called for some expected sysctls")
	}
}

func TestZpoolObjsetParsing(t *testing.T) {
	zpoolPaths, err := filepath.Glob("fixtures/proc/spl/kstat/zfs/*/objset-*")
	if err != nil {
		t.Fatal(err)
	}

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	for _, zpoolPath := range zpoolPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			t.Fatal(err)
		}

		err = c.parsePoolObjsetFile(file, zpoolPath, func(poolName string, datasetName string, s zfsSysctl, v uint64) {
			if s != zfsSysctl("kstat.zfs.misc.objset.writes") {
				return
			}

			handlerCalled = true

			if v != uint64(0) && v != uint64(4) && v != uint64(10) {
				t.Fatalf("Incorrect value parsed from procfs data %v", v)
			}

		})
		file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	if !handlerCalled {
		t.Fatal("Zpool parsing handler was not called for some expected sysctls")
	}
}

func TestAbdstatsParsing(t *testing.T) {
	abdstatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/abdstats")
	if err != nil {
		t.Fatal(err)
	}
	defer abdstatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(abdstatsFile, "abdstats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.abdstats.linear_data_size") {
			return
		}

		handlerCalled = true

		if v != uint64(223232) {
			t.Fatalf("Incorrect value parsed from procfs abdstats data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("ABDStats parsing handler was not called for some expected sysctls")
	}
}

func TestDbufstatsParsing(t *testing.T) {
	dbufstatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/dbufstats")
	if err != nil {
		t.Fatal(err)
	}
	defer dbufstatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(dbufstatsFile, "dbufstats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.dbufstats.hash_hits") {
			return
		}

		handlerCalled = true

		if v != uint64(108807) {
			t.Fatalf("Incorrect value parsed from procfs dbufstats data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("DbufStats parsing handler was not called for some expected sysctls")
	}
}

func TestDnodestatsParsing(t *testing.T) {
	dnodestatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/dnodestats")
	if err != nil {
		t.Fatal(err)
	}
	defer dnodestatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(dnodestatsFile, "dnodestats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.dnodestats.dnode_hold_alloc_hits") {
			return
		}

		handlerCalled = true

		if v != uint64(37617) {
			t.Fatalf("Incorrect value parsed from procfs dnodestats data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Dnodestats parsing handler was not called for some expected sysctls")
	}
}

func TestVdevMirrorstatsParsing(t *testing.T) {
	vdevMirrorStatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/vdev_mirror_stats")
	if err != nil {
		t.Fatal(err)
	}
	defer vdevMirrorStatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseProcfsFile(vdevMirrorStatsFile, "vdev_mirror_stats", func(s zfsSysctl, v uint64) {

		if s != zfsSysctl("kstat.zfs.misc.vdev_mirror_stats.preferred_not_found") {
			return
		}

		handlerCalled = true

		if v != uint64(94) {
			t.Fatalf("Incorrect value parsed from procfs vdev_mirror_stats data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("VdevMirrorStats parsing handler was not called for some expected sysctls")
	}
}

func TestPoolStateParsing(t *testing.T) {
	zpoolPaths, err := filepath.Glob("fixtures/proc/spl/kstat/zfs/*/state")
	if err != nil {
		t.Fatal(err)
	}

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	for _, zpoolPath := range zpoolPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			t.Fatal(err)
		}

		err = c.parsePoolStateFile(file, zpoolPath, func(poolName string, stateName string, isActive uint64) {
			handlerCalled = true

			if poolName == "pool1" {
				if isActive != uint64(1) && stateName == "online" {
					t.Fatalf("Incorrect parsed value for online state")
				}
				if isActive != uint64(0) && stateName != "online" {
					t.Fatalf("Incorrect parsed value for online state")
				}
			}
			if poolName == "poolz1" {
				if isActive != uint64(1) && stateName == "degraded" {
					t.Fatalf("Incorrect parsed value for degraded state")
				}
				if isActive != uint64(0) && stateName != "degraded" {
					t.Fatalf("Incorrect parsed value for degraded state")
				}
			}
			if poolName == "pool2" {
				if isActive != uint64(1) && stateName == "suspended" {
					t.Fatalf("Incorrect parsed value for suspended state")
				}
				if isActive != uint64(0) && stateName != "suspended" {
					t.Fatalf("Incorrect parsed value for suspended state")
				}
			}

		})
		file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	if !handlerCalled {
		t.Fatal("Zpool parsing handler was not called for some expected sysctls")
	}

}
