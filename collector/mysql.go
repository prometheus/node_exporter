// +build !nonetstat

package collector

import (
	"fmt"
	"strconv"
	"errors"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlStatsSubsystem = "mysql"
)

type mysqlStatCollector struct {
	config  Config
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories["mysql"] = NewMySQLStatCollector
}

func NewMySQLStatCollector(config Config) (Collector, error) {
	return &mysqlStatCollector{
		config:  config,
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *mysqlStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	// default
	mysqlConnString := "root:@/mysql"

	if c.config.Config["mysql_connection"] != "" {
		mysqlConnString = c.config.Config["mysql_connection"]
	}

	mysqlStats, err := getStats(mysqlConnString)
	if err != nil {
		return fmt.Errorf("couldn't get mysql stats: %s", err)
	}

	for key, value := range mysqlStats {
		if _, ok := c.metrics[key]; !ok {
			c.metrics[key] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: Namespace,
					Subsystem: mysqlStatsSubsystem,
					Name:      key,
					},
				)
			}

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value %s in mysqlStats: %s", value, err)
		}

		c.metrics[key].Set(v)
	}

	for _, m := range c.metrics {
		m.Collect(ch)
	}

	return err
}

// http://www.techrepublic.com/blog/linux-and-open-source/10-mysql-variables-that-you-should-monitor/
// http://blog.webyog.com/2012/09/03/top-10-things-to-monitor-on-your-mysql/
func getStats(mysqlConnString string) (map[string]string, error) {
	result := make(map[string]string)

	db, err := sql.Open("mysql", mysqlConnString)
	if err != nil {
		return nil, errors.New("error connecting to db")
	}
	defer db.Close()

	rows, err := db.Query("SHOW STATUS")
	if err != nil {
		return nil, errors.New("error running query")
	}
	defer rows.Close()

	var key, val []byte
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, errors.New("error getting row")
		}

		fmt.Printf("%s - %s\n", string(key), string(val))

		v := string(val)
		if _, err := strconv.Atoi(v); err == nil {
			result[string(key)] = v
		}

	}

	return result, nil
}
