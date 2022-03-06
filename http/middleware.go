// Copyright 2022 arcadium.dev <info@arcadium.dev>
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

package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"arcadium.dev/core/log"
)

var (
	concurrentRequests prometheus.Gauge
	requestCount       *prometheus.CounterVec
	requestSeconds     *prometheus.CounterVec
	requestTime        *prometheus.HistogramVec
)

func init() {
	concurrentRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_concurrent_requests",
		Help: "The number of concurrent http requests being processed",
	})

	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "Total number of requests by route",
	}, []string{"method", "path"})

	requestSeconds = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_seconds",
		Help: "Total amount of request time by route, in seconds",
	}, []string{"method", "path"})

	requestTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_seconds_histogram",
		Help: "The request time by route, in seconds.",
	}, []string{"method", "path"})
}

// TrackMetrics is middleware which provides tracking of key metrics for incoming
// HTTP requests.
func TrackMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track concurrent requests.
		concurrentRequests.Inc()
		defer concurrentRequests.Dec()

		// Obtain the template for the requested route.
		tmpl := ""
		if route := mux.CurrentRoute(r); route != nil {
			var err error
			if tmpl, err = route.GetPathTemplate(); err != nil {
				log.Error("msg", "failed to get path template", "route", route, "error", err.Error())
			}
		}

		// Time the handler.
		t := time.Now()

		// Delegate to the next handler in the middleware chain.
		next.ServeHTTP(w, r)

		// Track the request.
		s := float64(time.Since(t).Seconds())
		requestCount.WithLabelValues(r.Method, tmpl).Inc()
		requestSeconds.WithLabelValues(r.Method, tmpl).Add(s)
		requestTime.WithLabelValues(r.Method, tmpl).Observe(s)
	})
}
