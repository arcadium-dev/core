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

package middleware // import "arcadium.dev/core/middleware"

import (
	"net/http"
	"runtime"

	"github.com/rs/zerolog"
)

type (
	Recover struct {
		Logger *zerolog.Logger
	}
)

// Panics is middleware for recovering and reporting panics.
func (m Recover) Panics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Install the recovery and reporting function.
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				m.Logger.Error().Msg("recovering from a panic")

				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				m.Logger.Error().Msgf("stacktrace: %s", string(buf))
			}
		}()
		// Delegate to next handler in middleware chain.
		next.ServeHTTP(w, r)
	})
}
