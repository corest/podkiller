package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("Checheck /health endpoint by %s\n", r.RemoteAddr)
	data, _ := json.Marshal(healthCheckResponse{Status: "Pod killer is up and running"})
	writeJSONResponse(w, http.StatusOK, data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func healthHandler() {
	h := http.NewServeMux()
	h.HandleFunc("/health", health)
	go http.ListenAndServe(":9090", h)
}
