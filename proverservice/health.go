package main

import (
	"encoding/json"
	"net/http"
)

type status string

const (
	up   status = "UP"
	down status = "DOWN"
)

type HealthCheckResponse struct {
	Status status
}

func health(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HealthCheckResponse{Status: up})
}
