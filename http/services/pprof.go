// Copyright 2021-2024 arcadium.dev <info@arcadium.dev>
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

package services // import "arcadium.dev/core/http/services"

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

const (
	PProfRoute = "/debug/pprof/"
)

type (
	PProf struct{}
)

// Register sets up the http handler for this service with the given router.
func (p PProf) Register(router *mux.Router) {
	router.PathPrefix(PProfRoute).Handler(http.DefaultServeMux)
}

// Name returns the name of the service.
func (PProf) Name() string {
	return "pprof"
}

// Shutdown is a no-op since there no long running processes for this service.
func (PProf) Shutdown() {}
