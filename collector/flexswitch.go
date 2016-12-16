package collector

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	flexswitchPortStatsIgnoredPorts = flag.String(
		"collector.flexswitch.ignored-ports", "^$",
		"Regexp of port interfaces to ignore for flexswitch port collector.")
	flexswitchHost = flag.String(
		"collector.flexswitch.host", "localhost",
		"Hostname to use for REST query.")
	flexswitchProto = flag.String(
		"collector.flexswitch.proto", "http",
		"Protocol to use for REST query")
	flexswitchPort = flag.String(
		"collector.flexswitch.port", "8080",
		"Port to use for REST query.")
	procPortIntFieldSep = regexp.MustCompile("[ :] *")
)

type flexswitchStatsCollector struct {
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	metricDescs           map[string]*prometheus.Desc
}

type FlexPortDetail struct {
	IntfRef                     string  `json:"IntfRef"`
	IfIndex                     float64 `json:"IfIndex"`
	Name                        string  `json:"Name"`
	OperState                   string  `json:"OperState"`
	NumUpEvents                 float64 `json:"NumUpEvents"`
	LastUpEventTime             string  `json:"LastUpEventTime"`
	NumDownEvents               float64 `json:"NumDownEvents"`
	LastDownEventTim            string  `json:"LastDownEventTim"`
	Pvid                        float64 `json:"Pvid"`
	IfInOctets                  float64 `json:"IfInOctets"`
	IfInUcastPkts               float64 `json:"IfInUcastPkts"`
	IfInDiscards                float64 `json:"IfInDiscards"`
	IfInErrors                  float64 `json:"IfInErrors"`
	IfInUnknownProtos           float64 `json:"IfInUnknownProtos"`
	IfOutOctets                 float64 `json:"IfOutOctets"`
	IfOutUcastPkts              float64 `json:"IfOutUcastPkts"`
	IfOutDiscards               float64 `json:"IfOutDiscards"`
	IfOutErrors                 float64 `json:"IfOutErrors"`
	IfEtherUnderSizePktCnt      float64 `json:"IfEtherUnderSizePktCnt"`
	IfEtherOverSizePktCnt       float64 `json:"IfEtherOverSizePktCnt"`
	IfEtherFragments            float64 `json:"IfEtherFragments"`
	IfEtherCRCAlignError        float64 `json:"IfEtherCRCAlignError"`
	IfEtherJabber               float64 `json:"IfEtherJabber"`
	IfEtherPkts                 float64 `json:"IfEtherPkts"`
	IfEtherMCPkts               float64 `json:"IfEtherMCPkts"`
	IfEtherBcastPkts            float64 `json:"IfEtherBcastPkts"`
	IfEtherPkts64OrLessOctets   float64 `json:"IfEtherPkts64OrLessOctets"`
	IfEtherPkts65To127Octets    float64 `json:"IfEtherPkts65To127Octets"`
	IfEtherPkts128To255Octets   float64 `json:"IfEtherPkts128To255Octets"`
	IfEtherPkts256To511Octets   float64 `json:"IfEtherPkts256To511Octets"`
	IfEtherPkts512To1023Octets  float64 `json:"IfEtherPkts512To1023Octets"`
	IfEtherPkts1024To1518Octets float64 `json:"IfEtherPkts1024To1518Octets"`
	ErrDisableReason            string  `json:"ErrDisableReason"`
	PresentInHW                 string  `json:"PresentInHW"`
	ConfigMode                  string  `json:"ConfigMode"`
	PRBSRxErrCnt                float64 `json:"PRBSRxErrCnt"`
}

type FlexPort struct {
	ObjectId string         `"json:"ObjectId"`
	Object   FlexPortDetail `"json:Object"`
}

type FlexPortIndex struct {
	MoreExist     bool       `json"MoreExist"`
	ObjCount      float64    `json"ObjCount"`
	CurrentMarker float64    `json"CurrentMarker"`
	NextMarker    float64    `json"NextMarker"`
	Objects       []FlexPort `json"Objects"`
}

type FlexPortStat map[float64]FlexPort

func init() {
	Factories["flexswitch"] = NewPortStatsCollector
}

func NewPortStatsCollector() (Collector, error) {
	pattern := regexp.MustCompile(*flexswitchPortStatsIgnoredPorts)
	return &flexswitchStatsCollector{
		subsystem:             "network",
		ignoredDevicesPattern: pattern,
		metricDescs:           map[string]*prometheus.Desc{},
	}, nil
}

func getFlexswitchNetDevStats(ignore *regexp.Regexp) (map[string]map[string]string, error) {
	flexswitchPortsUrl := *flexswitchProto + "://" +
		*flexswitchHost + ":" + *flexswitchPort +
		"/public/v1/state/ports"
	resp, err := http.Get(flexswitchPortsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseFlexSwitchStats(resp, ignore)
}

func parseFlexSwitchStats(r *http.Response, ignore *regexp.Regexp) (map[string]map[string]string, error) {
	htmlBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid http body output:", err)
	}

	var jsonBody FlexPortIndex
	err = json.Unmarshal(htmlBody, &jsonBody)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json output:", err)
	}

	netDev := map[string]map[string]string{}

	for i := 0; i < int(jsonBody.ObjCount); i++ {
		Port := jsonBody.Objects[i].Object
		if ignore.MatchString(Port.Name) {
			log.Debugf("Ignoring device: %s", Port.Name)
			continue
		}
		netDev[Port.Name] = map[string]string{}
		// We should probably not be in the business of renaming the output
		// but I want them to be consistent with /proc/net/dev... least for
		// the ones that sync up
		netDev[Port.Name]["receive_bytes"] = strconv.FormatFloat(Port.IfInOctets, 'E', -1, 64)
		netDev[Port.Name]["transmit_bytes"] = strconv.FormatFloat(Port.IfOutOctets, 'E', -1, 64)

		netDev[Port.Name]["receive_packets"] = strconv.FormatFloat(Port.IfInUcastPkts, 'E', -1, 64)
		netDev[Port.Name]["transmit_packets"] = strconv.FormatFloat(Port.IfOutUcastPkts, 'E', -1, 64)

		netDev[Port.Name]["receive_errs"] = strconv.FormatFloat(Port.IfInErrors, 'E', -1, 64)
		netDev[Port.Name]["transmit_errs"] = strconv.FormatFloat(Port.IfOutErrors, 'E', -1, 64)

		netDev[Port.Name]["receive_drop"] = strconv.FormatFloat(Port.IfInDiscards, 'E', -1, 64)
		netDev[Port.Name]["transmit_drop"] = strconv.FormatFloat(Port.IfOutDiscards, 'E', -1, 64)

		netDev[Port.Name]["NumUpEvents"] = strconv.FormatFloat(Port.NumUpEvents, 'E', -1, 64)
		netDev[Port.Name]["NumDownEvents"] = strconv.FormatFloat(Port.NumDownEvents, 'E', -1, 64)
		netDev[Port.Name]["IfInUnknownProtos"] = strconv.FormatFloat(Port.IfInUnknownProtos, 'E', -1, 64)
		netDev[Port.Name]["IfEtherUnderSizePktCnt"] = strconv.FormatFloat(Port.IfEtherUnderSizePktCnt, 'E', -1, 64)
		netDev[Port.Name]["IfEtherOverSizePktCnt"] = strconv.FormatFloat(Port.IfEtherOverSizePktCnt, 'E', -1, 64)
		netDev[Port.Name]["IfEtherFragments"] = strconv.FormatFloat(Port.IfEtherFragments, 'E', -1, 64)
		netDev[Port.Name]["IfEtherCRCAlignError"] = strconv.FormatFloat(Port.IfEtherCRCAlignError, 'E', -1, 64)
		netDev[Port.Name]["IfEtherJabber"] = strconv.FormatFloat(Port.IfEtherJabber, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts"] = strconv.FormatFloat(Port.IfEtherPkts, 'E', -1, 64)
		netDev[Port.Name]["receive_multicast"] = strconv.FormatFloat(Port.IfEtherMCPkts, 'E', -1, 64)
		netDev[Port.Name]["IfEtherBcastPkts"] = strconv.FormatFloat(Port.IfEtherBcastPkts, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts64OrLessOctets"] = strconv.FormatFloat(Port.IfEtherPkts64OrLessOctets, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts65To127Octets"] = strconv.FormatFloat(Port.IfEtherPkts65To127Octets, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts128To255Octets"] = strconv.FormatFloat(Port.IfEtherPkts128To255Octets, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts256To511Octets"] = strconv.FormatFloat(Port.IfEtherPkts256To511Octets, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts512To1023Octets"] = strconv.FormatFloat(Port.IfEtherPkts512To1023Octets, 'E', -1, 64)
		netDev[Port.Name]["IfEtherPkts1024To1518Octets"] = strconv.FormatFloat(Port.IfEtherPkts1024To1518Octets, 'E', -1, 64)
	}

	return netDev, nil
}

func (c *flexswitchStatsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	netDev, err := getFlexswitchNetDevStats(c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get flexswitch port stats: %s", err)
	}
	for dev, devStats := range netDev {
		for key, value := range devStats {
			desc, ok := c.metricDescs[key]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, c.subsystem, key),
					fmt.Sprintf("flexswitch network device statistic %s.", key),
					[]string{"device"},
					nil,
				)
				c.metricDescs[key] = desc
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in port stats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, dev)
		}
	}
	return nil
}
