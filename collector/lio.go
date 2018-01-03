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
    lio_fileio_Subsystem = "lio_fileio"
    lio_iblock_Subsystem = "lio_iblock"
    lio_rbd_Subsystem    = "lio_rbd"
    lio_rdmcp_Subsystem  = "lio_rd_mcp"
)

// An lioCollector is a Collector which gathers iscsi RBD
// iops (iscsi commands) , Read in byte and Write in byte. 
// ( original reading sysfs is in MB )

type lioCollector struct {
    fs      sysfs.FS
    metrics *lioMetric
}

type lioMetric struct {
    lio_file_iops   *prometheus.Desc
    lio_file_read   *prometheus.Desc
    lio_file_write  *prometheus.Desc

    lio_block_iops  *prometheus.Desc
    lio_block_read  *prometheus.Desc
    lio_block_write *prometheus.Desc

    lio_rbd_iops    *prometheus.Desc
    lio_rbd_read    *prometheus.Desc
    lio_rbd_write   *prometheus.Desc

    lio_rdmcp_iops  *prometheus.Desc
    lio_rdmcp_read  *prometheus.Desc
    lio_rdmcp_write *prometheus.Desc
}

type graph_label struct {
    iqn     string
    tpgt    string
    lun     string
    store   string
    pool    string
    image   string
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

    metrics, _ := NewLioMetric()
    
    return &lioCollector{
        fs: fs,
        metrics: metrics }, nil
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

func NewLioMetric() (*lioMetric, error) {

    return &lioMetric{
        lio_file_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "iops"),
            "iSCSI FileIO backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
        ),
        lio_file_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "read"),
            "iSCSI FileIO backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
        ),
        lio_file_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "write"),
            "iSCSI FileIO backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "fileio", "object", "filename"}, nil,
        ),

        lio_block_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "iops"),
            "iSCSI IBlock backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
        ),
        lio_block_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "read"),
            "iSCSI IBlock backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
        ),
        lio_block_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "write"),
            "iSCSI IBlock backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "iblock", "object", "blockname"}, nil,
        ),

        lio_rbd_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "iops"),
            "iSCSI RBD backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
        ),
        lio_rbd_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "read"),
            "iSCSI RBD backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
        ),
        lio_rbd_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "write"),
            "iSCSI RBD backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "rbd", "pool", "image"}, nil,
        ),

        lio_rdmcp_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "iops"),
            "iSCSI Memory Copy RAMDisk backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "rd_mcp", "object"}, nil,
        ),
        lio_rdmcp_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "read"),
            "iSCSI Memory Copy RAMDisk backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "rd_mcp", "object" }, nil,
        ),
        lio_rdmcp_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "write"),
            "iSCSI Memory Copy RAMDisk backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "rd_mcp", "object" }, nil,
        ),
    }, nil
}

// /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*/ which link
// back to the following 
// /sys/kernel/config/target/core/{backstore_type}_{number}/{object_name}/

func (c *lioCollector) updateStat(ch chan<- prometheus.Metric, s *iscsi.Stats) error {

    log.Debugf("lio updateStat iscsi %s path", s.Name)
    tpgt_s := s.Tpgt
    for _, tpgt := range tpgt_s {
        tpgt_path := tpgt.Tpgt_path
        iscsi_enable := tpgt.Is_enable

        log.Debugf("lio: iscsi %s isEnable=%t", tpgt_path, iscsi_enable)
        // let's not putting more line into the graph with multiple
        // disable lun, it may create problem for bigger cluster
        if (iscsi_enable) {

            lun_s := tpgt.Luns
            for _, lun := range lun_s {
                backstore_type  := lun.Backstore
                object_name     := lun.Object_name
                type_number     := lun.Type_number

                // struct type graph_label { iqn, tpgt, lun, store, pool,  image}
                // label := graph_label {iqn, tpgt, lun, backstore_type, object_name, type_number}
                label := graph_label { s.Name, tpgt.Name, lun.Name, backstore_type, object_name, type_number}

                log.Debugf("lio: iqn=%s, tpgt=%s, lun=%s, type=%s, object=%s, type_number=%s", 
                s.Name, tpgt.Name, lun.Name, backstore_type, object_name, type_number) 

                switch { 
                    case strings.Compare(backstore_type, "fileio") == 0:
                        if err := c.updateFileIOStat(ch, label); err != nil {
                            return fmt.Errorf("failed fileio stat : %v", err)
                        }
                        break;
                    case strings.Compare(backstore_type, "iblock") == 0:
                        if err := c.updateIBlockStat(ch, label); err != nil {
                            return fmt.Errorf("failed iblock stat : %v", err)
                        }
                        break;
                    case strings.Compare(backstore_type, "rbd") == 0:
                        if err := c.updateRBDStat(ch, label); err != nil {
                            return fmt.Errorf("failed rbd stat : %v", err)
                        }
                        break;
                    case strings.Compare(backstore_type, "rd_mcp") == 0:
                        if err := c.updateRDMCPStat(ch, label); err != nil {
                            return fmt.Errorf("failed rdmcp stat : %v", err)
                        }
                        break;
                    default:
                        continue
                }
            }
        }
    }
    return nil
}

// /sys/kernel/config/target/core/fileio_{type_number}/{object}/
// udev_path has the file name 
func (c *lioCollector) updateFileIOStat(ch chan<- prometheus.Metric, label graph_label) error {

    fileio := new(iscsi.FILEIO)
    fileio, err := fileio.GetFileioUdev(label.image, label.pool)
    if err != nil {
        return err
    } 

    read_mb, write_mb, iops, err:= iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
    if err != nil {
        return err
    }
    log.Debugf("lio: Fileio Read int %d", read_mb) 
    f_read_mb  := float64(read_mb<<20) 
    log.Debugf("lio: Fileio Read float %f", f_read_mb) 

    log.Debugf("lio: Fileio Write int %d", write_mb) 
    f_write_mb := float64(write_mb<<20)
    log.Debugf("lio: Fileio Write int %f", f_write_mb) 

    log.Debugf("lio: Fileio OPS int %d", iops) 
    f_iops := float64(iops)
    log.Debugf("lio: Fileio OPS float %f", f_iops) 

    ch <- prometheus.MustNewConstMetric(c.metrics.lio_file_read,
    prometheus.CounterValue, f_read_mb, label.iqn, label.tpgt, label.lun, 
    fileio.Name, fileio.Object_name, fileio.Filename)

    ch <- prometheus.MustNewConstMetric(c.metrics.lio_file_write,
    prometheus.CounterValue, f_write_mb, label.iqn, label.tpgt, label.lun,
    fileio.Name, fileio.Object_name, fileio.Filename)

    ch <- prometheus.MustNewConstMetric(c.metrics.lio_file_iops,
    prometheus.CounterValue, f_iops, label.iqn, label.tpgt, label.lun,
    fileio.Name, fileio.Object_name, fileio.Filename)

    return nil
}

// /sys/kernel/config/target/core/iblock_{type_number}/{object}/
// udev_path has the file name 
func (c *lioCollector) updateIBlockStat(ch chan<- prometheus.Metric, label graph_label) error {

    iblock := new(iscsi.IBLOCK)
    iblock, err := iblock.GetIblockUdev(label.image, label.pool)
    if err != nil {
        return err
    }
    read_mb, write_mb, iops, err:= iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
    if err != nil {
        return err
    }
    log.Debugf("lio: IBlock Read int %d", read_mb) 
    f_read_mb  := float64(read_mb<<20) 
    log.Debugf("lio: IBlock Read float %f", f_read_mb) 

    log.Debugf("lio: IBlock Write int %d", write_mb) 
    f_write_mb := float64(write_mb<<20)
    log.Debugf("lio: IBlock Write int %f", f_write_mb) 

    log.Debugf("lio: IBlock OPS int %d", iops) 
    f_iops := float64(iops)
    log.Debugf("lio: IBlock OPS float %f", f_iops) 


    ch <- prometheus.MustNewConstMetric(c.metrics.lio_block_read,
    prometheus.CounterValue, f_read_mb, label.iqn, label.tpgt, label.lun, 
    iblock.Name, iblock.Object_name, iblock.Iblock)

    ch <- prometheus.MustNewConstMetric(c.metrics.lio_block_write,
    prometheus.CounterValue, f_write_mb, label.iqn, label.tpgt, label.lun,
    iblock.Name, iblock.Object_name, iblock.Iblock)

    ch <- prometheus.MustNewConstMetric(c.metrics.lio_block_iops,
    prometheus.CounterValue, f_iops, label.iqn, label.tpgt, label.lun,
    iblock.Name, iblock.Object_name, iblock.Iblock)

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

func (c *lioCollector) updateRBDStat(ch chan<- prometheus.Metric, label graph_label) error {

    rbd := new(iscsi.RBD)
    rbd, err := rbd.GetRBDMatch(label.image, label.pool)
    if err != nil {
        return err
    }
    if rbd != nil { 
        read_mb, write_mb, iops, err:= iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
        if err != nil {
            return err
        }
        log.Debugf("lio: RBD Read int %d", read_mb) 
        f_read_mb  := float64(read_mb<<20) 
        log.Debugf("lio: RBD Read float %f", f_read_mb) 

        log.Debugf("lio: RBD Write int %d", write_mb) 
        f_write_mb := float64(write_mb<<20)
        log.Debugf("lio: RBD Write int %f", f_write_mb) 

        log.Debugf("lio: RBD OPS int %d", iops) 
        f_iops := float64(iops)
        log.Debugf("lio: RBD OPS float %f", f_iops) 

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rbd_read,
        prometheus.CounterValue, f_read_mb, label.iqn, label.tpgt, label.lun,
        rbd.Name, rbd.Pool, rbd.Image)

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rbd_write,
        prometheus.CounterValue, f_write_mb, label.iqn, label.tpgt, label.lun,
        rbd.Name, rbd.Pool, rbd.Image)

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rbd_iops,
        prometheus.CounterValue, f_iops, label.iqn, label.tpgt, label.lun,
        rbd.Name, rbd.Pool, rbd.Image)
    }
    return nil
}

// /sys/kernel/config/target/core/rd_mcp_{type_number}/{object}/
// there won't be udev_path for ramdisk so not image name either 
func (c *lioCollector) updateRDMCPStat(ch chan<- prometheus.Metric, label graph_label) error {
    rd_mcp := new(iscsi.RDMCP)
    rd_mcp, err := rd_mcp.GetRDMCPPath(label.image, label.pool)
    if err != nil {
        return err
    }
    if rd_mcp != nil { 
        read_mb, write_mb, iops, err:= iscsi.ReadWriteOPS(label.iqn, label.tpgt, label.lun)
        if err != nil {
            return err
        }
        log.Debugf("lio: RBD Read int %d", read_mb) 
        f_read_mb  := float64(read_mb<<20) 
        log.Debugf("lio: RBD Read float %f", f_read_mb) 

        log.Debugf("lio: RBD Write int %d", write_mb) 
        f_write_mb := float64(write_mb<<20)
        log.Debugf("lio: RBD Write int %f", f_write_mb) 

        log.Debugf("lio: RBD OPS int %d", iops) 
        f_iops := float64(iops)
        log.Debugf("lio: RBD OPS float %f", f_iops) 

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rdmcp_read,
        prometheus.CounterValue, f_read_mb, label.iqn, label.tpgt, label.lun,
        rd_mcp.Name, rd_mcp.Object_name)

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rdmcp_write,
        prometheus.CounterValue, f_write_mb, label.iqn, label.tpgt, label.lun,
        rd_mcp.Name, rd_mcp.Object_name)

        ch <- prometheus.MustNewConstMetric(c.metrics.lio_rdmcp_iops,
        prometheus.CounterValue, f_iops, label.iqn, label.tpgt, label.lun,
        rd_mcp.Name, rd_mcp.Object_name)
    }
    return nil
}

