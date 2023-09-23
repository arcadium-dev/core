// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package provides a set of simple services that can be used the the http
// server.
package services // import "arcadium.dev/core/services"

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"arcadium.dev/core/build"
)

const (
	HealthRoute string = "/health"
)

type (
	// Health reports on the health of the service as a whole.
	Health struct {
		Start time.Time
		Info  build.Information
	}

	// HealthResponse returns the status of the server, including the
	// name, version, build details, and uptime.
	//
	// swagger:response HealthResponse
	HealthResponse struct {
		// swagger:allOf
		Data `json:"data"`
	}

	Data struct {
		// swagger:allOf
		Build   `json:"build"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Uptime  string `json:"uptime"`
	}

	Build struct {
		Branch string `json:"branch"`
		Commit string `json:"commit"`
		Date   string `json:"date"`
		Go     string `json:"go"`
	}
)

// Register sets up the http handler for this service with the given router.
func (h Health) Register(router *mux.Router) {
	// swagger:route GET /health Status health-id
	//
	// Reports the server health.
	//
	// Responses:
	//	200: HealthResponse
	r := router.PathPrefix(HealthRoute).Subrouter()
	r.HandleFunc("", h.get).Methods(http.MethodGet)
}

// Name returns the name of the service.
func (Health) Name() string {
	return "health"
}

// Shutdown is a no-op since there no long running processes for this service.
func (Health) Shutdown() {}

func (h Health) get(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{}
	resp.Data.Name = h.Info.Name
	resp.Data.Version = h.Info.Version
	resp.Data.Build.Branch = h.Info.Branch
	resp.Data.Build.Commit = h.Info.Commit
	resp.Data.Build.Date = h.Info.Date
	resp.Data.Build.Go = h.Info.Go
	resp.Data.Uptime = time.Since(h.Start).String()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}
