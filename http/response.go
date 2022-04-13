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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	cerrors "arcadium.dev/core/errors"
	"arcadium.dev/core/log"
)

// Response writes an http error responses to the http.ResponseWriter.
func Response(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	logger := log.LoggerFromContext(ctx)

	switch {
	case errors.Is(err, cerrors.ErrInvalidArgument):
		logger.Warn("reason", err.Error())
		response(ctx, w, http.StatusBadRequest, err)

	case errors.Is(err, cerrors.ErrNotFound):
		logger.Warn("reason", err.Error())
		response(ctx, w, http.StatusNotFound, err)

	case errors.Is(err, cerrors.ErrAlreadyExists):
		logger.Warn("reason", err.Error())
		response(ctx, w, http.StatusConflict, err)

	case errors.Is(err, cerrors.ErrNotImplemented):
		logger.Error("error", err.Error())
		response(ctx, w, http.StatusNotImplemented, err)

	default:
		logger.Error("error", err.Error())
		response(ctx, w, http.StatusInternalServerError, err)
	}
}

func response(ctx context.Context, w http.ResponseWriter, status int, e error) {
	err := ResponseError{Status: status}
	if e != nil {
		err.Detail = e.Error()
	}
	w.WriteHeader(err.Status)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	httpErrorCount.WithLabelValues(strconv.Itoa(err.Status)).Inc()

	resp := struct {
		Error ResponseError `json:"error,omitempty"`
	}{
		Error: err,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.LoggerFromContext(ctx).Error(
			"msg", "unable to write error response", "error", err.Error(),
		)
	}
}

type (
	// ResponseError provides additional information about problems encounted while
	// performing an operation. See: https://jsonapi.org/format/#error-objects
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
