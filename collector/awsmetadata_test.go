package collector

import (
	"testing"
)

func TestAwsMetadata(t *testing.T) {
	var wantevent = [2]map[string]string{}
	wantevent[0] = map[string]string{
		"NotBefore" : "1 May 2021 22:00:00 GMT", 
		"NotAfter" : "2 May 2021 00:00:00 GMT",
		"State" : "active",
	}
	wantevent[1] = map[string]string{
		"NotBefore" : "3 May 2021 22:00:00 GMT", 
		"NotAfter" : "4 May 2021 00:00:00 GMT",
		"State" : "active",
	}
	var wantmetrics = [2]map[string]float64{}
	wantmetrics[0] = map[string]float64{
		"notbefore" : 1619906400,
		"notafter" : 1619913600,
		"state" : 1,
	}
	wantmetrics[1] = map[string]float64{
		"notbefore" : 1620079200,
		"notafter" : 1620086400,
		"state" : 1,
	}

	events, err := parseAwsScheduledEvents(`[ {
  "NotBefore" : "1 May 2021 22:00:00 GMT",
  "Code" : "system-reboot",
  "Description" : "scheduled reboot do-not-complete",
  "EventId" : "instance-event-0000a0aa0aa0a0aaa",
  "NotAfter" : "2 May 2021 00:00:00 GMT",
  "State" : "active"
},
{
  "NotBefore" : "3 May 2021 22:00:00 GMT",
  "Code" : "system-reboot",
  "Description" : "scheduled reboot do-not-complete",
  "EventId" : "instance-event-0000b0bb0bb0b0bbb",
  "NotAfter" : "4 May 2021 00:00:00 GMT",
  "State" : "active"
} ]`)

	if err != nil {
		t.Fatal(err)
	}

	for i, event := range events {
		if event.State != wantevent[i]["State"] {
			t.Fatalf("want event State %s, got %s", wantevent[i]["State"], event.State)
		}
		if event.NotBefore != wantevent[i]["NotBefore"] {
			t.Fatalf("want event NotBefore %s, got %s", wantevent[i]["NotBefore"], event.NotBefore)
		}
		if event.NotAfter != wantevent[i]["NotAfter"] {
			t.Fatalf("want event NotAfter %s, got %s", wantevent[i]["NotAfter"], event.NotAfter)
		}

		eventMetrics, err := parseAwsScheduledEventMetrics(event)

		if err != nil {
			t.Fatal(err)
		}

		if eventMetrics[0] != wantmetrics[i]["state"] {
			t.Fatalf("want metric state %f, got %f", wantmetrics[i]["state"], eventMetrics[0])
		}
		if eventMetrics[1] != wantmetrics[i]["notbefore"] {
			t.Fatalf("want metric notbefore %f, got %f", wantmetrics[i]["notbefore"], eventMetrics[0])
		}
		if eventMetrics[2] != wantmetrics[i]["notafter"] {
			t.Fatalf("want metric notafter %f, got %f", wantmetrics[i]["notafter"], eventMetrics[0])
		}
	}

}
