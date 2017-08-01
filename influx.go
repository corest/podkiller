package main

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	podkillerDb       = "podkiller"
	podkillerDbTag    = "podkiller"
	podkillerDbSeries = "doomedpod"
)

// InfluxManager influxdb manager
type InfluxManager struct {
	client     client.Client
	doomedPods []*DoomedPod
}

func getInfluxClient(config *Config) client.Client {
	var clnt client.Client

	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", config.Influx.Hostname, config.Influx.Port),
	})
	if err != nil {
		log.Fatal(err)
	}

	return clnt
}

func (inflx *InfluxManager) initDB() {
	q := fmt.Sprintf("CREATE DATABASE %s", podkillerDb)
	_, err := inflx.queryDB(q)
	if err != nil {
		log.Fatalf("Failed to initialize database %s %v", podkillerDb, err)
	}
}

// queryDB convenience function to query the database
func (inflx *InfluxManager) queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: podkillerDb,
	}

	if response, err := inflx.client.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func (inflx *InfluxManager) addDoomedPod(pod *DoomedPod) {
	inflx.doomedPods = append(inflx.doomedPods, pod)
}

func (inflx *InfluxManager) writePodNecrology() error {

	if len(inflx.doomedPods) > 0 {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  podkillerDb,
			Precision: "ms",
		})
		if err != nil {
			log.Fatal(err)
		}

		var fields map[string]interface{}
		tags := map[string]string{"service": podkillerDbTag}

		for i, pod := range inflx.doomedPods {
			fields = map[string]interface{}{
				"podname":       pod.name,
				"namespace":     pod.namespace,
				"condemnded_at": pod.condemnedAt,
			}
			pt, err := client.NewPoint(podkillerDbSeries, tags, fields,
				time.Now().Add(time.Millisecond*time.Duration(i)))
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)
		}

		// Write the batch
		if err := inflx.client.Write(bp); err != nil {
			log.Fatal(err)
		}
		log.Printf("Submitted %d records into series '%s'", len(inflx.doomedPods), podkillerDbSeries)
		inflx.doomedPods = nil
	}
	return nil
}
