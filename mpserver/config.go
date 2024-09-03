// Copyright 2024 arcadium.dev <info@arcadium.dev>
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

package mpserver // import "arcadium.dev/mpserver"

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Config holds the configuration information for a multiprotocol server.
	Config struct {
		logLevel        string
		shutdownTimeout time.Duration
	}
)

// NewConfig returns the configuration of a multiprotocol server.
func NewConfig(prefix ...string) (Config, error) {
	cfg := struct {
		LogLevel        string        `split_words:"true"`
		ShutdownTimeout time.Duration `split_words:"true"`
	}{
		LogLevel:        DefaultLogLevel,
		ShutdownTimeout: DefaultShutdownTimeout,
	}

	pfix := ""
	if len(prefix) == 1 {
		pfix = prefix[0]
	}
	if err := envconfig.Process(pfix, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	c := Config{
		logLevel:        strings.TrimSpace(cfg.LogLevel),
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	return c, nil
}

// LogLevel returns the logging level. The default level is "info".
func (c Config) LogLevel() string {
	return c.logLevel
}

// ShutdownTimeout returns the server shutdown timeout.
func (c Config) ShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}
