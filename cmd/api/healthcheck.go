package main

import (
	"encoding/json"
	"net/http"

	"github.com/jahidhimon/greenlight.git/internal/greenlog"
)

type healthStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (a *application) Healthcheckhandler(w http.ResponseWriter, r *http.Request) {
	status := healthStatus{
		Status:      "available",
		Environment: a.config.env,
		Version:     version,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds information that we want to send as response
	data := map[string]string{
		"satus":       "available",
		"environment": a.config.env,
		"version":     version,
	}
	// Pass the map to the json.Marshal method. It returns a byte slice
	// containing encoded json
	err := a.writeJson(w, http.StatusOK, envelop{"health_status": data}, nil)
	// If there was a error, we log it and send the client a generic error message
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	greenlog.Logreq(r, "HealthCheck")
}
