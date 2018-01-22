// Copyright 2017 Alex Lau (AvengerMoJo) <alau@suse.com>
//
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
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs/iscsi"
	"github.com/prometheus/procfs/sysfs"
)

const (
	lioFileioSubsystem = "lio_fileio"
	lioIblockSubsystem = "lio_iblock"
	lioRbdSubsystem    = "lio_rbd"
	lioRdmcpSubsystem  = "lio_rdmcp"
)

// An lioCollector is a Collector which gathers iscsi RBD
// iops (iscsi commands) , Read in byte and Write in byte.
// ( original reading sysfs is in MB )

type lioCollector struct {
	fs      sysfs.FS
	metrics *lioMetric
}

type lioMetric struct {
	lioFileIops  *prometheus.Desc
	lioFileRead  *prometheus.Desc
	lioFileWrite *prometheus.Desc

	lioBlockIops  *prometheus.Desc
	lioBlockRead  *prometheus.Desc
	lioBlockWrite *prometheus.Desc

	lioRbdIops  *prometheus.Desc
	lioRbdRead  *prometheus.Desc
	lioRbdWrite *prometheus.Desc

	lioRdmcpIops  *prometheus.Desc
	lioRdmcpRead  *prometheus.Desc
	lioRdmcpWrite *prometheus.Desc
}

type graphLabel struct {
	iqn   string
	tpgt  string
	lun   string
	store string
	pool  string
	image string
}

func init() {
	registerCollector("iscsi", defaultEnabled, NewLioCollector)
}

// NewLioCollector returns a new Collector with iscsi statistics.
func NewLioCollector() (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %v", err)
	}

	metrics, _ := newLioMetric()

	return &lioCollector{
		fs:      fs,
		metrics: metrics}, nil
}

// Update implement the lioCollector.
func (c *lioCollector) Update(ch chan<- prometheus.Metric) error {

	stats, err := c.fs.ISCSIStats()
	log.Debugf("lio: Update lioCollector")
	if err != nil {
		return fmt.Errorf("lio: failed to update iscsi stat : %v", err)
	}
	for _, s := range stats {
		c.updateStat(ch, s)
	}
	return nil
}

//newLioMetric create the LIO metric data structure to return for node_exporter
func newLioMetric() (*lioMetric, error) {

	return &lioMetric{
		lioFileIops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioFileioSubsystem, "iops"),
			"iSCSI FileIO backstore transport operations.",
			[]string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
		),
		lioFileRead: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioFileioSubsystem, "read"),
			"iSCSI FileIO backstore Read in byte.",
			[]string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
		),
		lioFileWrite: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioFileioSubsystem, "write"),
			"iSCSI FileIO backstore Write in byte.",
			[]string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
		),

		lioBlockIops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioIblockSubsystem, "iops"),
			"iSCSI IBlock backstore transport operations.",
			[]string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
		),
		lioBlockRead: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioIblockSubsystem, "read"),
			"iSCSI IBlock backstore Read in byte.",
			[]string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
		),
		lioBlockWrite: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioIblockSubsystem, "write"),
			"iSCSI IBlock backstore Write in byte.",
			[]string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
		),

		lioRbdIops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRbdSubsystem, "iops"),
			"iSCSI RBD backstore transport operations.",
			[]string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
		),
		lioRbdRead: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRbdSubsystem, "read"),
			"iSCSI RBD backstore Read in byte.",
			[]string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
		),
		lioRbdWrite: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRbdSubsystem, "write"),
			"iSCSI RBD backstore Write in byte.",
			[]string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
		),

		lioRdmcpIops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRdmcpSubsystem, "iops"),
			"iSCSI Memory Copy RAMDisk backstore transport operations.",
			[]string{"iqn", "tpgt", "lun", "rdmcp", "object"}, nil,
		),
		lioRdmcpRead: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRdmcpSubsystem, "read"),
			"iSCSI Memory Copy RAMDisk backstore Read in byte.",
			[]string{"iqn", "tpgt", "lun", "rdmcp", "object"}, nil,
		),
		lioRdmcpWrite: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, lioRdmcpSubsystem, "write"),
			"iSCSI Memory Copy RAMDisk backstore Write in byte.",
			[]string{"iqn", "tpgt", "lun", "rdmcp", "object"}, nil,
		),
	}, nil
}

// /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*/ which link
// back to the following
// /sys/kernel/config/target/core/{backstoreType}_{number}/{objectName}/

func (c *lioCollector) updateStat(ch chan<- prometheus.Metric, s *iscsi.Stats) error {

	log.Debugf("lio updateStat iscsi %s path", s.Name)
	tpgtS := s.Tpgt
	for _, tpgt := range tpgtS {
		tpgtPath := tpgt.TpgtPath
		iscsiEnable := tpgt.IsEnable

		log.Debugf("lio: iscsi %s isEnable=%t", tpgtPath, iscsiEnable)
		// let's not putting more line into the graph with multiple
		// disable lun, it may create problem for bigger cluster
		if iscsiEnable {

			lunS := tpgt.Luns
			for _, lun := range lunS {
				backstoreType := lun.Backstore
				objectName := lun.ObjectName
				typeNumber := lun.TypeNumber

				// struct type graphLabel { iqn, tpgt, lun, store, pool,  image}
				// label := graphLabel {iqn, tpgt, lun, backstoreType, objectName, typeNumber}
				label := graphLabel{s.Name, tpgt.Name, lun.Name, backstoreType, objectName, typeNumber}

				log.Debugf("lio: iqn=%s, tpgt=%s, lun=%s, type=%s, object=%s, typeNumber=%s",
					s.Name, tpgt.Name, lun.Name, backstoreType, objectName, typeNumber)

				switch {
				case strings.Compare(backstoreType, "fileio") == 0:
					if err := c.updateFileIOStat(ch, label); err != nil {
						return fmt.Errorf("failed fileio stat : %v", err)
					}
					break
				case strings.Compare(backstoreType, "iblock") == 0:
					if err := c.updateIBlockStat(ch, label); err != nil {
						return fmt.Errorf("failed iblock stat : %v", err)
					}
					break
				case strings.Compare(backstoreType, "rbd") == 0:
					if err := c.updateRBDStat(ch, label); err != nil {
						return fmt.Errorf("failed rbd stat : %v", err)
					}
					break
				case strings.Compare(backstoreType, "rdmcp") == 0:
					if err := c.updateRDMCPStat(ch, label); err != nil {
						return fmt.Errorf("failed rdmcp stat : %v", err)
					}
					break
				default:
					continue
				}
			}
		}
	}
	return nil
}

// /sys/kernel/config/target/core/fileio_{typeNumber}/{object}/
// udev_path has the file name
func (c *lioCollector) updateFileIOStat(ch chan<- prometheus.Metric, label graphLabel) error {

	fileio := new(iscsi.FILEIO)
	fileio, err := fileio.GetFileioUdev(label.image, label.pool)
	if err != nil {
		return err
	}

	readMB, writeMB, iops, err := iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
	if err != nil {
		return err
	}
	log.Debugf("lio: Fileio Read int %d", readMB)
	fReadMB := float64(readMB << 20)
	log.Debugf("lio: Fileio Read float %f", fReadMB)

	log.Debugf("lio: Fileio Write int %d", writeMB)
	fWriteMB := float64(writeMB << 20)
	log.Debugf("lio: Fileio Write int %f", fWriteMB)

	log.Debugf("lio: Fileio OPS int %d", iops)
	fIops := float64(iops)
	log.Debugf("lio: Fileio OPS float %f", fIops)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioFileRead,
		prometheus.CounterValue, fReadMB, label.iqn, label.tpgt, label.lun,
		fileio.Name, fileio.ObjectName, fileio.Filename)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioFileWrite,
		prometheus.CounterValue, fWriteMB, label.iqn, label.tpgt, label.lun,
		fileio.Name, fileio.ObjectName, fileio.Filename)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioFileIops,
		prometheus.CounterValue, fIops, label.iqn, label.tpgt, label.lun,
		fileio.Name, fileio.ObjectName, fileio.Filename)

	return nil
}

// /sys/kernel/config/target/core/iblock_{typeNumber}/{object}/
// udev_path has the file name
func (c *lioCollector) updateIBlockStat(ch chan<- prometheus.Metric, label graphLabel) error {

	iblock := new(iscsi.IBLOCK)
	iblock, err := iblock.GetIblockUdev(label.image, label.pool)
	if err != nil {
		return err
	}
	readMB, writeMB, iops, err := iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
	if err != nil {
		return err
	}
	log.Debugf("lio: IBlock Read int %d", readMB)
	fReadMB := float64(readMB << 20)
	log.Debugf("lio: IBlock Read float %f", fReadMB)

	log.Debugf("lio: IBlock Write int %d", writeMB)
	fWriteMB := float64(writeMB << 20)
	log.Debugf("lio: IBlock Write int %f", fWriteMB)

	log.Debugf("lio: IBlock OPS int %d", iops)
	fIops := float64(iops)
	log.Debugf("lio: IBlock OPS float %f", fIops)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioBlockRead,
		prometheus.CounterValue, fReadMB, label.iqn, label.tpgt, label.lun,
		iblock.Name, iblock.ObjectName, iblock.Iblock)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioBlockWrite,
		prometheus.CounterValue, fWriteMB, label.iqn, label.tpgt, label.lun,
		iblock.Name, iblock.ObjectName, iblock.Iblock)

	ch <- prometheus.MustNewConstMetric(c.metrics.lioBlockIops,
		prometheus.CounterValue, fIops, label.iqn, label.tpgt, label.lun,
		iblock.Name, iblock.ObjectName, iblock.Iblock)

	return nil
}

// First using the rbd device label to create all the state place holder,
// Base on the following:
// /sys/devices/rbd/{} [0-9]* as rbd{X}
// pool  = '/sys/devices/rbd/{X}/pool'
// image = '/sys/devices/rbd/{X}/name'
//
// Then we loop though the iscsi target and match the link with the above
// rbd info /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*/{symblink}
//
// The link location look something like as following
// /sys/kernel/config/target/core/rbd_{X}/{pool}-{images}/
//
// the rbd_{X} / {pool}-{image} should match the following

func (c *lioCollector) updateRBDStat(ch chan<- prometheus.Metric, label graphLabel) error {

	rbd := new(iscsi.RBD)
	rbd, err := rbd.GetRBDMatch(label.image, label.pool)
	if err != nil {
		return err
	}
	if rbd != nil {
		readMB, writeMB, iops, err := iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
		if err != nil {
			return err
		}
		log.Debugf("lio: RBD Read int %d", readMB)
		fReadMB := float64(readMB << 20)
		log.Debugf("lio: RBD Read float %f", fReadMB)

		log.Debugf("lio: RBD Write int %d", writeMB)
		fWriteMB := float64(writeMB << 20)
		log.Debugf("lio: RBD Write int %f", fWriteMB)

		log.Debugf("lio: RBD OPS int %d", iops)
		fIops := float64(iops)
		log.Debugf("lio: RBD OPS float %f", fIops)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRbdRead,
			prometheus.CounterValue, fReadMB, label.iqn, label.tpgt, label.lun,
			rbd.Name, rbd.Pool, rbd.Image)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRbdWrite,
			prometheus.CounterValue, fWriteMB, label.iqn, label.tpgt, label.lun,
			rbd.Name, rbd.Pool, rbd.Image)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRbdIops,
			prometheus.CounterValue, fIops, label.iqn, label.tpgt, label.lun,
			rbd.Name, rbd.Pool, rbd.Image)
	}
	return nil
}

// /sys/kernel/config/target/core/rdmcp_{typeNumber}/{object}/
// there won't be udev_path for ramdisk so not image name either
func (c *lioCollector) updateRDMCPStat(ch chan<- prometheus.Metric, label graphLabel) error {
	rdmcp := new(iscsi.RDMCP)
	rdmcp, err := rdmcp.GetRDMCPPath(label.image, label.pool)
	if err != nil {
		return err
	}
	if rdmcp != nil {
		readMB, writeMB, iops, err := iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
		if err != nil {
			return err
		}
		log.Debugf("lio: RDMCP Read int %d", readMB)
		fReadMB := float64(readMB << 20)
		log.Debugf("lio: RDMCP Read float %f", fReadMB)

		log.Debugf("lio: RDMCP Write int %d", writeMB)
		fWriteMB := float64(writeMB << 20)
		log.Debugf("lio: RDMCP Write int %f", fWriteMB)

		log.Debugf("lio: RDMCP OPS int %d", iops)
		fIops := float64(iops)
		log.Debugf("lio: RDMCP OPS float %f", fIops)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRdmcpRead,
			prometheus.CounterValue, fReadMB, label.iqn, label.tpgt, label.lun,
			rdmcp.Name, rdmcp.ObjectName)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRdmcpWrite,
			prometheus.CounterValue, fWriteMB, label.iqn, label.tpgt, label.lun,
			rdmcp.Name, rdmcp.ObjectName)

		ch <- prometheus.MustNewConstMetric(c.metrics.lioRdmcpIops,
			prometheus.CounterValue, fIops, label.iqn, label.tpgt, label.lun,
			rdmcp.Name, rdmcp.ObjectName)
	}
	return nil
}
