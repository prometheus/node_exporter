// Types for unmarshalling gmond's XML output.
//
// Not used elements in gmond's XML output are commented.
// In case you want to use them, please change the names so that one
// can understand without needing to know what the acronym stands for.
package ganglia

import "encoding/xml"

type ExtraElement struct {
	Name string `xml:"NAME,attr"`
	Val  string `xml:"VAL,attr"`
}

type ExtraData struct {
	ExtraElements []ExtraElement `xml:"EXTRA_ELEMENT"`
}

type Metric struct {
	Name  string  `xml:"NAME,attr"`
	Value float64 `xml:"VAL,attr"`
	/*
		Unit      string    `xml:"UNITS,attr"`
		Slope     string    `xml:"SLOPE,attr"`
		Tn        int       `xml:"TN,attr"`
		Tmax      int       `xml:"TMAX,attr"`
		Dmax      int       `xml:"DMAX,attr"`
	*/
	ExtraData ExtraData `xml:"EXTRA_DATA"`
}

type Host struct {
	Name string `xml:"NAME,attr"`
	/*
		Ip           string `xml:"IP,attr"`
		Tags         string `xml:"TAGS,attr"`
		Reported     int    `xml:"REPORTED,attr"`
		Tn           int    `xml:"TN,attr"`
		Tmax         int    `xml:"TMAX,attr"`
		Dmax         int    `xml:"DMAX,attr"`
		Location     string `xml:"LOCATION,attr"`
		GmondStarted int    `xml:"GMOND_STARTED",attr"`
	*/
	Metrics []Metric `xml:"METRIC"`
}

type Cluster struct {
	Name string `xml:"NAME,attr"`
	/*
		Owner     string `xml:"OWNER,attr"`
		LatLong   string `xml:"LATLONG,attr"`
		Url       string `xml:"URL,attr"`
		Localtime int    `xml:"LOCALTIME,attr"`
	*/
	Hosts []Host `xml:"HOST"`
}

type Ganglia struct {
	XMLNAME  xml.Name  `xml:"GANGLIA_XML"`
	Clusters []Cluster `xml:"CLUSTER"`
}
