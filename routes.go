package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type router struct {
	influxmanager *influxManager
}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func (rtr *router) health(w http.ResponseWriter, r *http.Request) {
	log.Printf("Check /health endpoint by %s\n", r.RemoteAddr)
	data, _ := json.Marshal(healthCheckResponse{Status: "Pod killer is up and running"})
	writeJSONResponse(w, http.StatusOK, data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (rtr *router) metrics(w http.ResponseWriter, r *http.Request) {
	log.Printf("Get /metrics endpoint by %s\n", r.RemoteAddr)
	responseCode := http.StatusOK
	var responseData []byte

	q := fmt.Sprintf("SELECT * FROM %s", podkillerDbSeries)
	res, err := rtr.influxmanager.queryDB(q)
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve data from database. %v", err)
		log.Printf(msg)
		responseData, _ = json.Marshal(errorResponse{Error: msg})
		responseCode = http.StatusInternalServerError
	} else {

		var mResponse metricsResponse
		if len(res) != 0 && len(res[0].Series) > 0 {
			series := res[0].Series[0]
			mResponse.Name = series.Name
			mResponse.Columns = series.Columns
			mResponse.Values = series.Values

			responseData, _ = json.Marshal(mResponse)
		} else {
			responseData, _ = json.Marshal(make(map[string]string))
		}
	}

	writeJSONResponse(w, responseCode, responseData)
}

type metricsResponse struct {
	Name    string          `json:"name"`
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
}

func routesHandler(config *killerConfig, influxmanager *influxManager) {
	h := http.NewServeMux()
	r := &router{influxmanager: influxmanager}
	h.HandleFunc("/health", r.health)
	h.HandleFunc("/metrics", r.metrics)
	go http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(config.General.Port)), h)
}
