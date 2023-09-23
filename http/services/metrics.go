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

package services // import "arcadium.dev/core/services"

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	MetricsRoute = "/metrics"
)

type (
	// Metrics that reports the metrics of the service.
	Metrics struct{}
)

// Register sets up the http handler for this service with the given router.
func (Metrics) Register(router *mux.Router) {
	// swagger:route GET /metrics Status metrics-id
	//
	// Used to scrape server metrics.
	//
	// Responses:
	//	200: EmptyResponse
	router.Handle(MetricsRoute, promhttp.Handler()).Methods(http.MethodGet)
}

// Name returns the name of the service.
func (Metrics) Name() string { return "metrics" }

// Shutdown is a no-op since there are no long running processes.
func (Metrics) Shutdown() {}
