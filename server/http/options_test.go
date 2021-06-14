// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
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
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	mocklog "arcadium.dev/core/log/mock"
)

func TestHTTPServerWithTLS(t *testing.T) {
	s := &Server{
		server: &http.Server{},
	}
	cfg := &tls.Config{}
	WithTLS(cfg).apply(s)

	if s.server.TLSConfig == nil {
		t.Error("failed to set TLSConfig")
	}
}

func TestHTTPServerWithLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocklog.NewMockLogger(ctrl)

	s := &Server{}

	WithLogger(mockLogger).apply(s)

	if s.logger != mockLogger {
		t.Error("failed to set logger")
	}
}
