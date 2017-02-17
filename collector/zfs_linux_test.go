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
	err = c.parseProcfsFile(arcstatsFile, "arcstats", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.arcstats.hits") {
			return
		}

		handlerCalled = true

		if v != int(8772612) {
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
	err = c.parseProcfsFile(zfetchstatsFile, "zfetchstats", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.zfetchstats.hits") {
			return
		}

		handlerCalled = true

		if v != int(7067992) {
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
	err = c.parseProcfsFile(zilFile, "zil", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.zil.zil_commit_count") {
			return
		}

		handlerCalled = true

		if v != int(10) {
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
	err = c.parseProcfsFile(vdevCacheStatsFile, "vdev_cache_stats", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.vdev_cache_stats.delegations") {
			return
		}

		handlerCalled = true

		if v != int(40) {
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
	err = c.parseProcfsFile(xuioStatsFile, "xuio_stats", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.xuio_stats.onloan_read_buf") {
			return
		}

		handlerCalled = true

		if v != int(32) {
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
	err = c.parseProcfsFile(fmFile, "fm", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.fm.erpt-dropped") {
			return
		}

		handlerCalled = true

		if v != int(18) {
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
	err = c.parseProcfsFile(dmuTxFile, "dmu_tx", func(s zfsSysctl, v int) {

		if s != zfsSysctl("kstat.zfs.misc.dmu_tx.dmu_tx_assigned") {
			return
		}

		handlerCalled = true

		if v != int(3532844) {
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

		err = c.parsePoolProcfsFile(file, zpoolPath, func(poolName string, s zfsSysctl, v int) {
			if s != zfsSysctl("kstat.zfs.misc.io.nread") {
				return
			}

			handlerCalled = true

			if v != int(1884160) && v != int(2826240) {
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
