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
// WITHOUT WARRANTIES OR CONDITIONS OF ANt KINr, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mockserver "arcadium.dev/core/server/mock"
)

func sharedNewTest(t *testing.T, setup func(*mockserver.MockConfig), check func(*Server, error), opts ...Option) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	setup(mockConfig)
	check(New(mockConfig, opts...))
}

func TestGRPCServerNewSecure01(t *testing.T) {
	sharedNewTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Cert().Return("../test_data/cert.pem")
			mockConfig.EXPECT().Key().Return("")
		},
		func(s *Server, err error) {
			expectedErr := "A certificate must be configured for TLS, or the WithInsecure option must be given to run without TLS."
			if s != nil || err == nil {
				t.Error("New failed")
			}
			if err.Error() != expectedErr {
				t.Errorf("incorrect error, expected: %s, actual: %s", expectedErr, err.Error())
			}
		},
	)
}

func TestGRPCServerNewSecure02(t *testing.T) {
	sharedNewTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Cert().Return("")
		},
		func(s *Server, err error) {
			expectedErr := "A certificate must be configured for TLS, or the WithInsecure option must be given to run without TLS."
			if s != nil || err == nil {
				t.Error("New failed")
			}
			if err.Error() != expectedErr {
				t.Errorf("incorrect error, expected: %s, actual: %s", expectedErr, err.Error())
			}
		},
	)
}

func TestGRPCServerNewTLS(t *testing.T) {
	sharedNewTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Cert().Return("../test_data/cert.pem").Times(2)
			mockConfig.EXPECT().Key().Return("../test_data/key.pem").Times(2)
			mockConfig.EXPECT().CACert().Return("../test_data/rootCA.pem")
			mockConfig.EXPECT().Addr().Return(":4201")
		},
		func(s *Server, err error) {
			if s == nil || err != nil {
				t.Error("New failed")
			}
		},
	)
}

func TestGRPCServerNewInsecure(t *testing.T) {
	sharedNewTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Addr().Return(":4201")
		},
		func(s *Server, err error) {
			if s == nil || err != nil {
				t.Error("New failed")
			}
		},
		WithInsecure(),
	)
}

func sharedServeTest(t *testing.T, setup func(*mockserver.MockConfig), check func(*Server, chan error)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mockserver.NewMockConfig(ctrl)
	setup(mockConfig)

	s, _ := New(mockConfig, WithInsecure())

	result := make(chan error)
	go s.Serve(result)

	check(s, result)
}

func TestGRPCServerServeBadAddr(t *testing.T) {
	sharedServeTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Addr().Return(":-1")
		},
		func(s *Server, result chan error) {
			err := <-result

			expectedErr := "Failed to listen on :-1: listen tcp: address -1: invalid port"
			if err.Error() != expectedErr {
				t.Errorf("incorrect error, expected: %s, actual: %s", expectedErr, err.Error())
			}
		},
	)
}

func TestGRPCServerServeStop(t *testing.T) {
	sharedServeTest(t,
		func(mockConfig *mockserver.MockConfig) {
			mockConfig.EXPECT().Addr().Return(":4201")
		},
		func(s *Server, result chan error) {
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
				t.Errorf("error should not have occurred, %s", err)
			}
		},
	)
}
