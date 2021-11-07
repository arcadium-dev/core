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

package mock // import "arcadium.dev/core/mock"

import (
	"arcadium.dev/core/log"
	"arcadium.dev/core/test"
)

type (
	// Logger is a mock implementation of the logger.
	Logger struct {
		Buffer *test.StringBuffer
		log.Logger
	}
)

// NewLogger returns a mock logger that uses a string buffer to gather logs.
func NewLogger(opts ...log.Option) (*Logger, error) {
	var err error

	logger := &Logger{
		Buffer: test.NewStringBuffer(),
	}
	opts = append(opts, log.WithOutput(logger.Buffer))
	logger.Logger, err = log.New(opts...)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
