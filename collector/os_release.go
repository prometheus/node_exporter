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

//go:build !noosrelease && !aix
// +build !noosrelease,!aix

package collector

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	envparse "github.com/hashicorp/go-envparse"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	etcOSRelease       = "/etc/os-release"
	usrLibOSRelease    = "/usr/lib/os-release"
	systemVersionPlist = "/System/Library/CoreServices/SystemVersion.plist"
)

var (
	versionRegex       = regexp.MustCompile(`^[0-9]+\.?[0-9]*`)
	nixOSCurrentSystem = "/run/current-system"
	nixOSBootedSystem  = "/run/booted-system"
	nixOSstoreDB       = "/nix/var/nix/db/db.sqlite"
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
	ImageID         string // NixOS
	ImageVersion    string
	SupportEnd      string

	ImageTime int // NixOS

	BootedImageTime int    // NixOS
	BootedImageID   string // NixOS
}

type osReleaseCollector struct {
	infoDesc           *prometheus.Desc
	logger             *slog.Logger
	os                 *osRelease
	osMutex            sync.RWMutex
	osReleaseFilenames []string // all os-release file names to check
	version            float64
	versionDesc        *prometheus.Desc
	supportEnd         time.Time
	supportEndDesc     *prometheus.Desc
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
func NewOSCollector(logger *slog.Logger) (Collector, error) {
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
		supportEndDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "os", "support_end_timestamp_seconds"),
			"Metric containing the end-of-life date timestamp of the OS.",
			nil, nil,
		),
	}, nil
}

func parseOSRelease(r io.Reader) (*osRelease, error) {
	env, err := envparse.Parse(r)
	result := osRelease{
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
		SupportEnd:      env["SUPPORT_END"],
	}

	if env["NAME"] == "NixOS" {
		nixCurrentPath := getRealPathFromLink(nixOSCurrentSystem)
		nixBootedPath := getRealPathFromLink(nixOSBootedSystem)

		result.ImageTime = getNixOSRegistrationTimeOfPath(nixCurrentPath)
		result.BootedImageTime = getNixOSRegistrationTimeOfPath(nixBootedPath)

		result.ImageID = getNixOSHashFromPath(nixOSCurrentSystem)
		result.BootedImageID = getNixOSHashFromPath(nixOSBootedSystem)
	}

	return &result, err
}

func getNixOSRegistrationTimeOfPath(path string) int {
	var registrationTime int
	uri := fmt.Sprintf("file:%s?immutable=true", nixOSstoreDB)

	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return registrationTime
	}
	defer db.Close()

	query := fmt.Sprintf("select registrationTime from ValidPaths where path = '%s'", path)
	rows, err := db.Query(query)
	if err != nil {
		return registrationTime
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&registrationTime)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = rows.Err()
	if err != nil {
		return registrationTime
	}

	return registrationTime
}

func getRealPathFromLink(path string) string {
	symLink, _ := filepath.EvalSymlinks(path)
	return symLink
}

func getNixOSHashFromPath(path string) string {
	symLink, _ := filepath.EvalSymlinks(path)
	symLinkHashLong := strings.Split(symLink, "/")[3]
	symLinkHashShort := strings.Split(symLinkHashLong, "-")[0]

	return string(symLinkHashShort)
}

func (c *osReleaseCollector) UpdateStruct(path string) error {
	releaseFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer releaseFile.Close()

	// Acquire a lock to update the osReleaseCollector struct.
	c.osMutex.Lock()
	defer c.osMutex.Unlock()

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

	if c.os.SupportEnd != "" {
		c.supportEnd, err = time.Parse(time.DateOnly, c.os.SupportEnd)

		if err != nil {
			return err
		}
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
				c.logger.Debug("no os-release file found", "files", strings.Join(c.osReleaseFilenames, ","))
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

	if c.os.SupportEnd != "" {
		ch <- prometheus.MustNewConstMetric(c.supportEndDesc, prometheus.GaugeValue, float64(c.supportEnd.Unix()))
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
