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

type influxManager struct {
	client     client.Client
	doomedPods []*doomedPod
}

func getInfluxClient(config *killerConfig) client.Client {
	var clnt client.Client

	clnt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", config.Influx.Hostname, config.Influx.Port),
	})
	if err != nil {
		log.Fatal(err)
	}

	return clnt
}

func (inflx *influxManager) initDB() {
	q := fmt.Sprintf("CREATE DATABASE %s", podkillerDb)
	_, err := inflx.queryDB(q)
	if err != nil {
		log.Fatalf("Failed to initialize database %s %v", podkillerDb, err)
	}
}

// queryDB convenience function to query the database
func (inflx *influxManager) queryDB(cmd string) (res []client.Result, err error) {
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

func (inflx *influxManager) addDoomedPod(pod *doomedPod) {
	inflx.doomedPods = append(inflx.doomedPods, pod)
}

func (inflx *influxManager) writePodNecrology() error {

	if len(inflx.doomedPods) > 0 {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  podkillerDb,
			Precision: "s",
		})
		if err != nil {
			log.Fatal(err)
		}

		var fields map[string]interface{}
		tags := map[string]string{"service": podkillerDbTag}

		for _, pod := range inflx.doomedPods {
			fields = map[string]interface{}{
				"name":          pod.name,
				"namespace":     pod.namespace,
				"condemnded_at": pod.condemnedAt,
			}
			pt, err := client.NewPoint(podkillerDbSeries, tags, fields, time.Now())
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
