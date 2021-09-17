package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type scheduledEvent struct {
  Code         string  `json:"code"`
  State        string  `json:"state"`
  Description  string  `json:"description"`
  EventID      string  `json:"eventid"`
  NotBefore    string  `json:"notbefore"`
  NotAfter     string  `json:"notafter"`
}

type awsmetadataCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("awsmetadata", defaultEnabled, NewAwsmetadataCollector)
}

// NewAwsmetadataCollector returns new Collector exposing AWS instance metadata stats
func NewAwsmetadataCollector(logger log.Logger) (Collector, error) {
	return &awsmetadataCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(namespace+"_state", "state of scheduled event", []string{"code", "state",}, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_notbefore", "earliest start time of scheduled event", []string{"code", "state",}, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_notafter", "latest start time of scheduled event", []string{"code", "state",}, nil), prometheus.GaugeValue},
		},
		logger: logger,
		}, nil
}

func (c *awsmetadataCollector) Update(ch chan<- prometheus.Metric) error {
	metricSets, err := c.getAwsMetadata()
	if err != nil {
		return fmt.Errorf("couldn't get scheduled events from instance metadata: %w", err)
	}

	for i, metrics := range metricSets {
		// TODO: start here -- need to setup metric Descs and push to channel
		for j, metric := range metrics {
			if j == 0 {
				level.Debug(c.logger).Log("msg", "return aws event", "index", i, "state", metric)
			} else if j == 1 {
				level.Debug(c.logger).Log("msg", "return aws event", "index", i, "notbefore", metric)
			} else {
				level.Debug(c.logger).Log("msg", "return aws event", "index", i, "notafter", metric)
			}

			ch <- c.metric[i].mustNewConstMetric(metric)
		}
	}

	return nil
}

// TODO: should i generalize this more?
// TODO: compare to https://github.com/aws/aws-node-termination-handler/blob/8eceda9337/pkg/ec2metadata/ec2metadata.go#L137
// return instance metadata collected through AWS IMDS
func (c *awsmetadataCollector) getAwsMetadata() ([][3]float64, error) {
	metrics := [][3]float64{}
	eventsMetadata, err := c.getAwsScheduledEvents()
	if err != nil {
		return nil, err
	}

	events, err := parseAwsScheduledEvents(eventsMetadata)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		eventMetrics, err := parseAwsScheduledEventMetrics(event)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, eventMetrics)
	}

	return metrics, nil
}

// get scheduled events via instance metadata
func (c *awsmetadataCollector) getAwsScheduledEvents() (string, error) {
	mdURL := "http://169.254.169.254/latest/meta-data/events/maintenance/scheduled"

	resp, err := http.Get(mdURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	mdEvents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(mdEvents), nil
}

// takes an array of json objects as a string and returns populated structs
func parseAwsScheduledEvents(data string) ([]scheduledEvent, error) {
	res := []scheduledEvent{}
	json.Unmarshal([]byte(data), &res)

	return res, nil
}

// returns metrics in the order {active, notbefore, notafter}
func parseAwsScheduledEventMetrics(event scheduledEvent) ([3]float64, error) {
	var metrics [3]float64
	tformat := "_2 Jan 2006 15:04:05 GMT"

	if event.State == "active" {
		metrics[0] = 1
	} else {
		metrics[0] = 0
	}

	nb, err := time.Parse(tformat, event.NotBefore)
	if err != nil {
		return [3]float64{0, 0, 0}, err
	}
	na, err := time.Parse(tformat, event.NotAfter)
	if err != nil {
		return [3]float64{0, 0, 0}, err
	}

	metrics[1] = float64(nb.Unix())
	metrics[2] = float64(na.Unix())

	return metrics, nil
}
