// Copyright 2021 The Prometheus Authors
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
	"encoding/xml"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	envparse "github.com/hashicorp/go-envparse"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	etcOSRelease       = "/etc/os-release"
	usrLibOSRelease    = "/usr/lib/os-release"
	systemVersionPlist = "/System/Library/CoreServices/SystemVersion.plist"
)

var (
	versionRegex = regexp.MustCompile(`^[0-9]+\.?[0-9]*`)
)

type osRelease struct {
	Name            string
	ID              string
	IDLike          string
	PrettyName      string
	Variant         string
	VariantID       string
	Version         string
	VersionID       string
	VersionCodename string
	BuildID         string
	ImageID         string
	ImageVersion    string
}

type osReleaseCollector struct {
	infoDesc           *prometheus.Desc
	logger             log.Logger
	os                 *osRelease
	osFilename         string    // file name of cached release information
	osMtime            time.Time // mtime of cached release file
	osMutex            sync.RWMutex
	osReleaseFilenames []string // all os-release file names to check
	version            float64
	versionDesc        *prometheus.Desc
}

type Plist struct {
	Dict Dict `xml:"dict"`
}

type Dict struct {
	Key    []string `xml:"key"`
	String []string `xml:"string"`
}

func init() {
	registerCollector("os", defaultEnabled, NewOSCollector)
}

// NewOSCollector returns a new Collector exposing os-release information.
func NewOSCollector(logger log.Logger) (Collector, error) {
	return &osReleaseCollector{
		logger: logger,
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "os", "info"),
			"A metric with a constant '1' value labeled by build_id, id, id_like, image_id, image_version, "+
				"name, pretty_name, variant, variant_id, version, version_codename, version_id.",
			[]string{"build_id", "id", "id_like", "image_id", "image_version", "name", "pretty_name",
				"variant", "variant_id", "version", "version_codename", "version_id"}, nil,
		),
		osReleaseFilenames: []string{etcOSRelease, usrLibOSRelease, systemVersionPlist},
		versionDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "os", "version"),
			"Metric containing the major.minor part of the OS version.",
			[]string{"id", "id_like", "name"}, nil,
		),
	}, nil
}

func parseOSRelease(r io.Reader) (*osRelease, error) {
	env, err := envparse.Parse(r)
	return &osRelease{
		Name:            env["NAME"],
		ID:              env["ID"],
		IDLike:          env["ID_LIKE"],
		PrettyName:      env["PRETTY_NAME"],
		Variant:         env["VARIANT"],
		VariantID:       env["VARIANT_ID"],
		Version:         env["VERSION"],
		VersionID:       env["VERSION_ID"],
		VersionCodename: env["VERSION_CODENAME"],
		BuildID:         env["BUILD_ID"],
		ImageID:         env["IMAGE_ID"],
		ImageVersion:    env["IMAGE_VERSION"],
	}, err
}

func (c *osReleaseCollector) UpdateStruct(path string) error {
	releaseFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	stat, err := releaseFile.Stat()
	if err != nil {
		return err
	}

	t := stat.ModTime()
	c.osMutex.RLock()
	upToDate := path == c.osFilename && t == c.osMtime
	c.osMutex.RUnlock()
	if upToDate {
		// osReleaseCollector struct is already up-to-date.
		return nil
	}

	// Acquire a lock to update the osReleaseCollector struct.
	c.osMutex.Lock()
	defer c.osMutex.Unlock()

	level.Debug(c.logger).Log("msg", "file modification time has changed",
		"file", path, "old_value", c.osMtime, "new_value", t)
	c.osFilename = path
	c.osMtime = t
	//  SystemVersion.plist is xml file with MacOs version info
	if strings.Contains(releaseFile.Name(), "SystemVersion.plist") {
		c.os, err = getMacosProductVersion(releaseFile.Name())
		if err != nil {
			return err
		}
	} else {
		c.os, err = parseOSRelease(releaseFile)
		if err != nil {
			return err
		}
	}

	majorMinor := versionRegex.FindString(c.os.VersionID)
	if majorMinor != "" {
		c.version, err = strconv.ParseFloat(majorMinor, 64)
		if err != nil {
			return err
		}
	} else {
		c.version = 0
	}
	return nil
}

func (c *osReleaseCollector) Update(ch chan<- prometheus.Metric) error {
	for i, path := range c.osReleaseFilenames {
		err := c.UpdateStruct(*rootfsPath + path)
		if err == nil {
			break
		}
		if errors.Is(err, os.ErrNotExist) {
			if i >= (len(c.osReleaseFilenames) - 1) {
				level.Debug(c.logger).Log("msg", "no os-release file found", "files", strings.Join(c.osReleaseFilenames, ","))
				return ErrNoData
			}
			continue
		}
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0,
		c.os.BuildID, c.os.ID, c.os.IDLike, c.os.ImageID, c.os.ImageVersion, c.os.Name, c.os.PrettyName,
		c.os.Variant, c.os.VariantID, c.os.Version, c.os.VersionCodename, c.os.VersionID)
	if c.version > 0 {
		ch <- prometheus.MustNewConstMetric(c.versionDesc, prometheus.GaugeValue, c.version,
			c.os.ID, c.os.IDLike, c.os.Name)
	}
	return nil
}

func getMacosProductVersion(filename string) (*osRelease, error) {
	f, _ := os.Open(filename)
	bytePlist, _ := io.ReadAll(f)
	f.Close()

	var plist Plist
	err := xml.Unmarshal(bytePlist, &plist)
	if err != nil {
		return &osRelease{}, err
	}

	var osVersionID, osVersionName, osBuildID string
	if len(plist.Dict.Key) > 0 {
		for index, value := range plist.Dict.Key {
			switch value {
			case "ProductVersion":
				osVersionID = plist.Dict.String[index]
			case "ProductName":
				osVersionName = plist.Dict.String[index]
			case "ProductBuildVersion":
				osBuildID = plist.Dict.String[index]
			}
		}
	}
	return &osRelease{
		Name:      osVersionName,
		Version:   osVersionID,
		VersionID: osVersionID,
		BuildID:   osBuildID,
	}, nil
}
