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

package grpc

//go:generate mockgen -package mockgrpc -destination ./mock/service.go . Service

import (
	"google.golang.org/grpc"
)

// Service is an abstraction of a single gRPC service.
type Service interface {
	// Register will register this service with the given gRPC server.
	Register(server *grpc.Server)

	// LogFields provides a set of log fields details.
	LogFields() map[string]interface{}
}
