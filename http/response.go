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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"arcadium.dev/core/log"
)

// Response writes http error responses to the http.ResponseWriter. The first
// error's status will be used to write the http response header.
func Response(w http.ResponseWriter, errs ...ResponseError) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for i := 0; i < len(errs); i++ {
		if errs[i].Status < 100 || 999 < errs[i].Status {
			errs[i].Status = http.StatusInternalServerError
		}
		if i == 0 {
			w.WriteHeader(errs[i].Status)
		}
		httpErrorCount.WithLabelValues(strconv.Itoa(errs[i].Status)).Inc()
	}

	if len(errs) > 0 {
		resp := struct {
			Errors []ResponseError `json:"errors,omitempty"`
		}{
			Errors: errs,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("msg", "unable to write error response", "error", err.Error())
		}
	}
}

// BadRequestError creates a ResponseError with the status of
// http.StatusBadRequest.
func BadRequestError(err error) ResponseError {
	return wrap(http.StatusBadRequest, err)
}

// NotFoundError creates a ResponseError with the status of
// http.StatusNotFound.
func NotFoundError(err error) ResponseError {
	return wrap(http.StatusNotFound, err)
}

// ConflictError creates a ResponseError with the status of
// http.StatusConflict.
func ConflictError(err error) ResponseError {
	return wrap(http.StatusConflict, err)
}

// InternalServerError creates a ResponseError with the status of
// http.StatusInternalServerError.
func InternalServerError(err error) ResponseError {
	return wrap(http.StatusInternalServerError, err)
}

// NotImplementedError creates a ResponseError with the status of
// http.StatusNotImplemented.
func NotImplementedError(err error) ResponseError {
	return wrap(http.StatusNotImplemented, err)
}

type (
	// ResponseError provides additional informatio about problems encounted
	// while performing an operation.
	//
	// See: https://jsonapi.org/format/#error-objects
	ResponseError struct {
		// Status is the http status code applicable to this problem.
		Status int `json:"status,omitempty"`
		// Detail is a human-readable explanation specific to this occurrence of
		// the problem.
		Detail string `json:"detail,omitempty"`
	}
)

// Error translates the error to a string.
func (e ResponseError) Error() string {
	return fmt.Sprintf("status=%d, detail=%q", e.Status, e.Detail)
}

var (
	httpErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_error_count",
		Help: "Total number of http errors by error status.",
	}, []string{"status"})
)

func wrap(status int, err error) ResponseError {
	rerr := ResponseError{Status: status}
	if err != nil {
		rerr.Detail = err.Error()
	}
	return rerr
}
