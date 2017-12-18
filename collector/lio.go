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
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "strconv"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/common/log"
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
    active  string
    store   string
    pool    string
    image   string
}

func init() {
    registerCollector("iscsi", defaultEnabled, NewLioCollector)
}

// NewLioCollector returns a new Collector with iscsi statistics.
func NewLioCollector() (Collector, error) {
    return &lioCollector{
        lio_file_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "iops"),
            "iSCSI FileIO backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "active", "fileio", "object", "filename"}, nil,
        ),
        lio_file_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "read"),
            "iSCSI FileIO backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "active", "fileio", "object", "filename"}, nil,
        ),
        lio_file_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_fileio_Subsystem, "write"),
            "iSCSI FileIO backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "active", "fileio", "object", "filename"}, nil,
        ),

        lio_block_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "iops"),
            "iSCSI IBlock backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "active", "iblock", "object", "blockname"}, nil,
        ),
        lio_block_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "read"),
            "iSCSI IBlock backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "active", "iblock", "object", "blockname"}, nil,
        ),
        lio_block_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_iblock_Subsystem, "write"),
            "iSCSI IBlock backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "active", "iblock", "object", "blockname"}, nil,
        ),

        lio_rbd_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "iops"),
            "iSCSI RBD backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "active", "rbd", "pool", "image"}, nil,
        ),
        lio_rbd_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "read"),
            "iSCSI RBD backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "active", "rbd", "pool", "image"}, nil,
        ),
        lio_rbd_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rbd_Subsystem, "write"),
            "iSCSI RBD backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "active", "rbd", "pool", "image"}, nil,
        ),

        lio_rdmcp_iops: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "iops"),
            "iSCSI Memory Copy RAMDisk backstore transport operations.",
            []string{"iqn", "tpgt", "lun", "active", "rd_mcp", "object", "size"}, nil,
        ),
        lio_rdmcp_read: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "read"),
            "iSCSI Memory Copy RAMDisk backstore Read in byte.",
            []string{"iqn", "tpgt", "lun", "active", "rd_mcp", "object", "size"}, nil,
        ),
        lio_rdmcp_write: prometheus.NewDesc(
            prometheus.BuildFQName(namespace, lio_rdmcp_Subsystem, "write"),
            "iSCSI Memory Copy RAMDisk backstore Write in byte.",
            []string{"iqn", "tpgt", "lun", "active", "rd_mcp", "object", "size"}, nil,
        ),
    }, nil
}

// Update implement the lioCollector.
func (c *lioCollector) Update(ch chan<- prometheus.Metric) error {

    log.Debugf("lioCollector updateStat\n")
    if err := c.updateStat(ch); err != nil {
        return fmt.Errorf("failed to update iscsi stat : %v", err)
    }
    return nil
}

// /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*/ which link
// back to the following 
// /sys/kernel/config/target/core/{backstore_type}_{number}/{object_name}/

func (c *lioCollector) updateStat(ch chan<- prometheus.Metric) error {
    iqn_s, _:=  getIQN() 

    for _, iqn_path :=range iqn_s {
        tpgt_s, _:= getTpgt(iqn_path) 

        for _, tpgt_path:= range tpgt_s { 
            iscsi_enable := isTpgtEnable(tpgt_path)

            log.Debugf("iscsi %s isEnable=%t\n", tpgt_path, iscsi_enable)

            // let's not putting more line into the graph with multiple
            // disable lun, it may create problem for bigger cluster

            if (iscsi_enable) {
                lun_s, _:= getLun(tpgt_path) 
                for _, lun_path := range lun_s {
                    backstore_type, object_name, type_number, err := getLunLinkTarget(lun_path)

                    if err != nil {
                        continue 
                    }
                    _, iqn := filepath.Split(iqn_path) 
                    _, tpgt:= filepath.Split(tpgt_path) 
                    _, lun := filepath.Split(lun_path) 

                    // struct type graph_label { iqn, tpgt, lun, active, store, pool,  image}
                    label := graph_label {iqn, tpgt, lun, "enable", backstore_type, object_name, type_number}

                    log.Debugf("iqn=%s, tpgt=%s, lun=%s, type=%s, object=%s, type_number=%s\n", 
                    iqn, tpgt, lun, backstore_type, object_name, type_number) 

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
                                return fmt.Errorf("failed rbd stat : %v", err)
                            }
                            break;
                        default:
                            continue
                    }
                }
            }
        }
    }
    return nil
}

// /sys/kernel/config/target/core/fileio_{type_number}/{object}/
// udev_path has the file name 
func (c *lioCollector) updateFileIOStat(ch chan<- prometheus.Metric, label graph_label) error {
    fileio_with_number := "fileio_" + label.image
    udev_path := filepath.Join("/sys/kernel/config/target/core/",
    fileio_with_number, label.pool, "udev_path")

    if _, err := os.Stat(udev_path); os.IsNotExist(err) {
        return fmt.Errorf("fileio_%s is missing file name ...!", label.image)
    } else {
        file_name, err := ioutil.ReadFile(udev_path)
        if err != nil {
            return fmt.Errorf("Cannot read file name from %s!", udev_path)
        } else { 
            label.image = strings.TrimSpace(string(file_name))
            log.Debugf("%s file name is using ->%s<- !",
            fileio_with_number, label.image)

            //  /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*

            log.Debugf("File %s\n", filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn, label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))
            read_mb, _:=readUintFromFile(
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))

            log.Debugf("Read %f\n", read_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_file_read,
            prometheus.CounterValue, float64(read_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            fileio_with_number, label.pool, label.image)

            write_mb, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/write_mbytes"))

            log.Debugf("Write %f\n", write_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_file_write,
            prometheus.CounterValue, float64(write_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            fileio_with_number, label.pool, label.image)

            iops, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/in_cmds"))

            log.Debugf("iops %f\n", iops) 

            ch <- prometheus.MustNewConstMetric(c.lio_file_iops,
            prometheus.CounterValue, float64(iops),
            label.iqn, label.tpgt, label.lun, "enable",
            fileio_with_number, label.pool, label.image)
        }
    }
    return nil
}

// /sys/kernel/config/target/core/iblock_{type_number}/{object}/
// udev_path has the file name 
func (c *lioCollector) updateIBlockStat(ch chan<- prometheus.Metric, label graph_label) error {
    iblock_with_number := "iblock_" + label.image
    udev_path := filepath.Join("/sys/kernel/config/target/core/",
    iblock_with_number, label.pool, "udev_path")

    if _, err := os.Stat(udev_path); os.IsNotExist(err) {
        return fmt.Errorf("iblock_%s is missing file name ...!", label.image)
    } else {
        block_name, err := ioutil.ReadFile(udev_path)
        if err != nil {
            return fmt.Errorf("Cannot read block name from %s!", udev_path)
        } else { 
            label.image = strings.TrimSpace(string(block_name))
            log.Debugf("%s block name is using ->%s<- !",
            iblock_with_number, label.image)

            //  /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*

            log.Debugf("File %s\n", filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn, label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))
            read_mb, _:=readUintFromFile(
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))

            log.Debugf("Read %f\n", read_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_block_read,
            prometheus.CounterValue, float64(read_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            iblock_with_number, label.pool, label.image)

            write_mb, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/write_mbytes"))

            log.Debugf("Write %f\n", write_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_block_write,
            prometheus.CounterValue, float64(write_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            iblock_with_number, label.pool, label.image)

            iops, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/in_cmds"))

            log.Debugf("iops %f\n", iops) 

            ch <- prometheus.MustNewConstMetric(c.lio_block_iops,
            prometheus.CounterValue, float64(iops),
            label.iqn, label.tpgt, label.lun, "enable",
            iblock_with_number, label.pool, label.image)
        }
    }
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

    rbds, err := filepath.Glob(sysFilePath("devices/rbd/[0-9]*"))
    if err != nil {
        return err
    }

    // looping all rbd{X} with image and pool
    for rbd_number , rbd_path := range rbds {

        log.Debugf("RBD path %s, device rbd%d", rbd_path, rbd_number)
        var pool string = ""
        var image string = ""

        log.Debugf("rbd_%d\n pool --->%s<--- \nimage --->%s<---", rbd_number, pool, image)

        if _, err := os.Stat(filepath.Join(rbd_path, "pool")); os.IsNotExist(err) {
            fmt.Errorf("rbd%d is missing pool name ...!", rbd_number)
            continue 
        } else {
            p_name , err := ioutil.ReadFile(filepath.Join(rbd_path, "pool"))
            if err != nil {
                fmt.Errorf("Cannot read pool name from rbd%d!", rbd_number)
                continue 
            } else {
                log.Debugf("rbd%d pool name is ->%s<- !", rbd_number, p_name)
                pool = strings.TrimSpace(string(p_name))
            }
        }
        if _, err := os.Stat(filepath.Join(rbd_path, "name")); os.IsNotExist(err) {
            fmt.Errorf("rbd%d is missing image name ...!", rbd_number)
            continue 
        } else {
            i_name, err := ioutil.ReadFile(filepath.Join(rbd_path, "name"))
            if err != nil {
                fmt.Errorf("Cannot read image name from rbd%d!", rbd_number)
                continue 
            } else { 
                log.Debugf("rbd%d image name is ->%s<- !", rbd_number, image)
                image = strings.TrimSpace(string(i_name))
            }
        }
        log.Debugf("rbd_%d\n pool --->%s<--- \nimage --->%s<---", rbd_number, pool, image)

        // find a match 
        // reminding label 
        // []string{"iqn", "tpgt", "lun", "active", "rbd", "pool", "image"}, nil,

        log.Debugf("Matching rbd%d, label->rbd%s", rbd_number, label.image )
        log.Debugf("Matching pool->%s, image->%s, label->%s", pool, image, label.pool)

        if matchRBD(fmt.Sprintf("%d", rbd_number), label.image) &&
        matchPoolImage(pool, image, label.pool) {
            log.Debugf("Match! rbd%d, pool->%s, image->%s", rbd_number, pool, image)

            //  /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*

            log.Debugf("File %s\n", filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn, label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))
            read_mb, _:=readUintFromFile(
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))

            log.Debugf("Read %f\n", read_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_rbd_read,
            prometheus.CounterValue, float64(read_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            "rbd" + label.image , pool, image)

            write_mb, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/write_mbytes"))

            log.Debugf("Write %f\n", write_mb) 

            ch <- prometheus.MustNewConstMetric(c.lio_rbd_write,
            prometheus.CounterValue, float64(write_mb<<10),
            label.iqn, label.tpgt, label.lun, "enable",
            "rbd" + label.image, pool, image)

            iops, _:=readUintFromFile( 
            filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
            label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/in_cmds"))

            log.Debugf("iops %f\n", iops) 

            ch <- prometheus.MustNewConstMetric(c.lio_rbd_iops,
            prometheus.CounterValue, float64(iops),
            label.iqn, label.tpgt, label.lun, "enable",
            "rbd" + label.image, pool, image)
        }
    }
    return nil
}

// /sys/kernel/config/target/core/rd_mcp_{type_number}/{object}/
// there won't be udev_path for ramdisk so not image name either 
func (c *lioCollector) updateRDMCPStat(ch chan<- prometheus.Metric, label graph_label) error {
    rd_mcp_with_number := "rd_mcp_" + label.image

    //  /sys/kernel/config/target/iscsi/iqn*/tpgt_*/lun/lun_*

    log.Debugf("File %s\n", filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn, label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))
    read_mb, _:=readUintFromFile(
    filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
    label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/read_mbytes"))

    log.Debugf("Read %f\n", read_mb) 

    ch <- prometheus.MustNewConstMetric(c.lio_rdmcp_read,
    prometheus.CounterValue, float64(read_mb<<10),
    label.iqn, label.tpgt, label.lun, "enable",
    rd_mcp_with_number, label.pool, "")

    write_mb, _:=readUintFromFile( 
    filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
    label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/write_mbytes"))

    log.Debugf("Write %f\n", write_mb) 

    ch <- prometheus.MustNewConstMetric(c.lio_rdmcp_write,
    prometheus.CounterValue, float64(write_mb<<10),
    label.iqn, label.tpgt, label.lun, "enable",
    rd_mcp_with_number, label.pool, "")

    iops, _:=readUintFromFile( 
    filepath.Join("/sys/kernel/config/target/iscsi/", label.iqn,
    label.tpgt, "lun", label.lun, "statistics/scsi_tgt_port/in_cmds"))

    log.Debugf("iops %f\n", iops) 

    ch <- prometheus.MustNewConstMetric(c.lio_rdmcp_iops,
    prometheus.CounterValue, float64(iops),
    label.iqn, label.tpgt, label.lun, "enable",
    rd_mcp_with_number, label.pool, "")

    return nil
}

func getIQN() (m []string, err error) { 
    matches, err := filepath.Glob(sysFilePath("kernel/config/target/iscsi/iqn*"))
    if err != nil {
        fmt.Errorf("getIQN error %v\n", err)
        return nil, nil
    }
    return matches, nil
}

func getTpgt(iqn_path string) (m []string, err error) {
    log.Debugf("getTpgt path %s\n", iqn_path)
    matches, err := filepath.Glob(filepath.Join(iqn_path, "tpgt*"))
    if err != nil {
        fmt.Errorf("getTPGT error %v\n", err )                                                                                                                           
        return nil, nil
    } 
    return matches, nil
}

func isTpgtEnable(tpgt_path string) (isEnable bool) { 
    isEnable = false
    tmp, err := ioutil.ReadFile(filepath.Join(tpgt_path, "enable"))
    log.Debugf("getTpgt enable %s\n", tmp)
    if err != nil {
        fmt.Errorf("isTpgtEnable error %v\n", err )                                                                                                                           
        return false
    } 
    tmp_num, err := strconv.Atoi(strings.TrimSpace(string(tmp)))
    log.Debugf("getTpgt enable number %d\n", tmp_num)
    if (tmp_num > 0) {
        isEnable = true
    }
    return isEnable
}

func getLun(tpgt_path string) (m []string, err error) { 
    log.Debugf("getLun path %s\n", tpgt_path)
    matches, err := filepath.Glob(filepath.Join(tpgt_path, "lun/lun*"))
    if err != nil {
        fmt.Errorf("getLun error  %v\n", err )
        return nil, nil
    } 
    return matches, nil
}

func getLunLinkTarget(lun_path string) ( backstore_type string, object_name string, type_number string, err error) { 
    files, err := ioutil.ReadDir(lun_path)
    if err != nil { 
        fmt.Errorf("getLunLinkTarget error  %v\n", err )
        return "", "", "", nil
    }
    for _, file := range files {
        log.Debugf("lun dir list file ->%s<-\n",file.Name())
        fileInfo, _:= os.Lstat(lun_path +  "/" + file.Name())
        if fileInfo.Mode() & os.ModeSymlink != 0 {
            target, err := os.Readlink( lun_path +  "/" + fileInfo.Name()) 
            if err != nil {
                fmt.Errorf("Readlink err %v\n", err)
                return "", "", "", nil
            }
            p1, object_name := filepath.Split(target)
            _, type_with_number := filepath.Split(filepath.Clean(p1))

            tmp := strings.Split(type_with_number, "_")
            backstore_type, type_number := tmp[0], tmp[1]
            if len(tmp) == 3 {
                backstore_type = fmt.Sprintf("%s_%s", tmp[0], tmp[1])
                type_number = tmp[2] 
            }
            log.Debugf("object_name->%s-<, type->%s, type_number->%s<- \n", object_name, backstore_type, type_number)
            return backstore_type, object_name, type_number, nil
        }
    }
    return "", "", "", errors.New("getLunLinkTarget: Lun Link does not exist")
}

func matchRBD(rbd_number string, rbd_name string) (isEqual bool) { 
    isEqual = false
    log.Debugf("compare rbd->%s<- with rbd->%s<- \n", rbd_number, rbd_name)
    if strings.Compare(rbd_name, rbd_number) == 0 { 
        isEqual = true
    }
    return isEqual 
}

func matchPoolImage(pool string, image string, match_pool_image string) (isEqual bool) { 
    isEqual = false
    var pool_image = fmt.Sprintf("%s-%s", pool, image)
    log.Debugf("compare pool_image->%s<- with pool_image->%s<- \n", pool_image, match_pool_image)
    if strings.Compare(pool_image, match_pool_image) == 0 { 
        isEqual = true
    }
    return isEqual 
}

