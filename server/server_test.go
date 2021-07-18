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

package server

import (
	"os"
	"testing"

	"arcadium.dev/core/errors"
)

func TestServerServeWithError(t *testing.T) {
	sig := make(chan os.Signal, 1)
	result := make(chan error, 1)

	err := serve(sig, result, mockServer{err: errors.New("server error")})
	if err == nil {
		t.Error("an error was expected")
	}
	if err.Error() != "server error" {
		t.Errorf("Incorrect error: expected 'server error', actual '%s'", err)
	}
}

func TestServerServeWithSuccess(t *testing.T) {
	sig := make(chan os.Signal, 1)
	result := make(chan error, 1)

	err := serve(sig, result, mockServer{})
	if err != nil {
		t.Error("an error was not expected")
	}
}

func TestServerServeWithSignal(t *testing.T) {
	sig := make(chan os.Signal, 1)
	result := make(chan error, 1)

	err := serve(sig, result, mockServer{sig: sig})
	if err != nil {
		t.Error("an error was not expected")
	}
}

type mockServer struct {
	err error
	sig chan os.Signal
}

func (s mockServer) Serve(result chan<- error) {
	if s.sig != nil {
		s.sig <- os.Interrupt
	}
	if s.err != nil {
		result <- s.err
	}
	close(result)
}

func (s mockServer) Stop() {}
