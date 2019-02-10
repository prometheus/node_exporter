package collector

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"os"
	"time"
	//  "strconv"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"path/filepath"
	"strings"
	//  "fmt"
)

const (
	certExpiration = "cert"
)

var (
	certLabelNames = []string{"main", "dns", "ip", "cert"}
	pathCert       = kingpin.Flag("collector.certificate.path", "Set proper key=value where value is a path to directory with certificates and key is only name without meaning.").Default("cert=/etc/ssl/certs").StringMap()
	pathExt        = kingpin.Flag("collector.certificate.ext", "Extension of certificates").Default("pem,crt").String()

	readsCert = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, certExpiration, "seconds_to_expire"),
		"The total seconds to expired certificate.",
		certLabelNames, nil,
	)
)

type certStatsCollector struct {
	certTime typedDesc
}

func init() {
	registerCollector("certificate", defaultDisabled, NewCertStatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewCertStatsCollector() (Collector, error) {
	return &certStatsCollector{
		certTime: typedDesc{readsCert, prometheus.CounterValue},
	}, nil
}

func (c *certStatsCollector) Update(ch chan<- prometheus.Metric) (err error) {

	var seconds float64
	for _, e := range checkExt() {
		var cert *x509.Certificate
		cert = decodeCert(e)

		//split full path of filename to get certificate name
		nameAfterSplit := strings.Split(e, "/")

		seconds = float64(convertToSeconds(cert.NotAfter))
		if err != nil {
			log.Infoln(err)
		}

		ch <- c.certTime.mustNewConstMetric(seconds, "true", "", "", strings.ToLower(nameAfterSplit[len(nameAfterSplit)-1]))
		for _, element := range cert.DNSNames {
			ch <- c.certTime.mustNewConstMetric(seconds, "false", element, "", strings.ToLower(nameAfterSplit[len(nameAfterSplit)-1]))
		}
		for _, element := range cert.IPAddresses {
			ch <- c.certTime.mustNewConstMetric(seconds, "false", "", element.String(), strings.ToLower(nameAfterSplit[len(nameAfterSplit)-1]))
		}
	}

	return nil
}

func convertToSeconds(expirationDate time.Time) float64 {
	now := time.Now()
	diff := expirationDate.Sub(now)

	//return strconv.Itoa(int(diff.Seconds()))
	return diff.Seconds()
}
func readCert(certFile string) string {
	dat, err := ioutil.ReadFile(certFile)
	//  check(err)
	if err != nil {
		log.Infoln(err)
		//    log.Fatal(err)
	}

	return string(dat)
}

func decodeCert(file string) *x509.Certificate {
	var certData = []byte(readCert(file))

	block, _ := pem.Decode(certData)
	//fmt.Print(block.Type)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Infoln("Failed to decode cert:", file)

		return nil
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Infoln("Failed to parse certificate: ", err)
		return nil
	}

	return cert
}

func checkExt() []string {

	var filesSSL []string
	for _, singlePath := range *pathCert {
		//    fmt.Print(singlePath + "\n")
		filepath.Walk(singlePath, func(path string, f os.FileInfo, _ error) error {
			if stringInArray(filepath.Ext(path), *pathExt) && f.Mode().IsRegular() {
				filesSSL = append(filesSSL, path)
			}

			return nil
		})
	}
	return filesSSL
}

func stringInArray(str string, stringToSplit string) bool {
	for _, v := range strings.Split(stringToSplit, ",") {
		if v == strings.Trim(str, ".") {
			return true
		}
	}
	return false
}
