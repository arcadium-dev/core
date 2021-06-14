// Copyright 2021 arcadium.dev <info@arcadium.dev>
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
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"arcadium.dev/core/server"
	mockserver "arcadium.dev/core/server/mock"
)

func TestHTTPServerNewWithTLS(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	mockConfig.EXPECT().Addr().Return(":8080")
	mockConfig.EXPECT().Cert().Return("../test_data/cert.pem")
	mockConfig.EXPECT().Key().Return("../test_data/key.pem")

	tlsConfig, tlsErr := server.CreateTLSConfig(mockConfig)
	if tlsErr != nil {
		t.Errorf("received tls error: %s", tlsErr)
	}

	s, err := New(mockConfig, WithTLS(tlsConfig))

	if s == nil {
		t.Error("failed to create http server")
	}
	if err != nil {
		t.Errorf("recieved unexpected error: %s", err)
	}
	if s.server.TLSConfig != tlsConfig {
		t.Error("failed to set tls config")
	}
}

func TestHttpServerNewWithoutTLS(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	mockConfig.EXPECT().Addr().Return(":8080")

	s, err := New(mockConfig)

	if s == nil {
		t.Error("failed to create http server")
	}
	if err != nil {
		t.Errorf("recieved unexpected error: %s", err)
	}
	if s.server.TLSConfig != nil {
		t.Error("tls config incorrectly set")
	}
}

func TestHttpServerHandleMux(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	mockConfig.EXPECT().Addr().Return(":8080")

	s, _ := New(mockConfig)

	mux := http.NewServeMux()

	s.Handle(mux)

	if s.server.Handler != mux {
		t.Error("failed to set handler")
	}
}

func TestHttpServerServeBadAddr(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	mockConfig.EXPECT().Addr().Return(":-1")

	s, _ := New(mockConfig)

	result := make(chan error)
	go s.Serve(result)
	err := <-result

	if err == nil {
		t.Error("expecting an error")
	}

	expectedErr := "Failed to listen on :-1: listen tcp: address -1: invalid port"
	if err.Error() != expectedErr {
		t.Errorf("expected error: %s, actual %s", expectedErr, err.Error())
	}
}

func TestHttpServerServeStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	mockConfig.EXPECT().Addr().Return(":8080")

	s, _ := New(mockConfig)

	result := make(chan error)
	go s.Serve(result)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
	defer cancel()

	var err error
	select {
	case err = <-result:
	case <-ctx.Done():
		s.Stop()
		err = <-result
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
