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

package grpc

/*
import (
	"testing"

	"github.com/golang/mock/gomock"

	// mocklog "arcadium.dev/core/log/mock"
	mocktrace "arcadium.dev/core/trace/mock"
)

func TestGRPCWithoutReflection(t *testing.T) {
	s := &Server{reflection: true}
	WithoutReflection().apply(s)

	if s.reflection != false || s.insecure != false || s.metrics != false ||
		s.logger != nil || s.tracer != nil {
		t.Errorf("WithoutReflection is broken")
	}
}

func TestGRPCWithInsecure(t *testing.T) {
	s := &Server{reflection: true}
	WithInsecure().apply(s)

	if s.reflection != true || s.insecure != true || s.metrics != false ||
		s.logger != nil || s.tracer != nil {
		t.Errorf("WithInsecure is broken")
	}
}

func TestGRPCWithMetrics(t *testing.T) {
	s := &Server{reflection: true}
	WithMetrics().apply(s)

	if s.reflection != true || s.insecure != false || s.metrics != true ||
		s.logger != nil || s.tracer != nil {
		t.Errorf("WithMetrics is broken")
	}
}

func TestGRPCWithLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocklog.NewMockLogger(ctrl)

	s := &Server{reflection: true}
	WithLogger(mockLogger).apply(s)

	if s.reflection != true || s.insecure != false || s.metrics != false ||
		s.logger != mockLogger || s.tracer != nil {
		t.Errorf("WithLogger is broken")
	}
}

func TestGRPCWithTracer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTracer := mocktrace.NewMockTracer(ctrl)

	s := &Server{reflection: true}
	WithTracer(mockTracer).apply(s)

	if s.reflection != true || s.insecure != false || s.metrics != false ||
		s.logger != nil || s.tracer != mockTracer {
		t.Errorf("WithTracer is broken")
	}
}
*/
