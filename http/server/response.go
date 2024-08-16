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

package server // import "arcadium.dev/core/http/server"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/rs/zerolog"

	"arcadium.dev/core/errors"
)

// Response writes an http error responses to the http.ResponseWriter.
func Response(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	respErr := ResponseError{
		Status: http.StatusInternalServerError,
		Detail: err.Error(),
	}

	var e errors.HTTPError
	if errors.As(err, &e) {
		respErr.Status = e.Status()
	} else {
		if c, ok := err.(interface{ Status() int }); ok {
			respErr.Status = c.Status()
		}
	}

	logger := zerolog.Ctx(ctx)
	switch {
	case errors.Is(err, context.Canceled):
		logger.Warn().Msg(err.Error())
		return
	case respErr.Status >= 500:
		logger.Err(err).Msg("")
	case respErr.Status >= 400:
		logger.Warn().Msgf("reason: %s", err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(respErr.Status)

	if encErr := json.NewEncoder(w).Encode(respErr); encErr != nil {
		zerolog.Ctx(ctx).Err(encErr).Msg("unable to write error response")
	}

	httpErrorCount.WithLabelValues(strconv.Itoa(respErr.Status)).Inc()
}

type (
	// ResponseError provides additional information about problems encounted while
	// performing an operation. See: https://jsonapi.org/format/#error-objects
	//
	// swagger:response ResponseError
	ResponseError struct {
		// Status is the http status code applicable to this problem.
		Status int `json:"status"`

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
