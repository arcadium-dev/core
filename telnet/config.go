// Copyright 2026 arcadium.dev <info@arcadium.dev>
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

package telnet // import "arcadium.dev/telnet"

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// ServerConfig holds the configuration information for the telnet server.
	ServerConfig struct {
		serverAddr string
	}
)

// NewServerConfig returns the configuration of telnet server.
func NewServerConfig(prefix ...string) (ServerConfig, error) {
	cfg := struct {
		ServerAddr string `required:"true" split_words:"true"`
	}{}

	pfix := ""
	if len(prefix) == 1 {
		pfix = prefix[0]
	}
	if err := envconfig.Process(pfix, &cfg); err != nil {
		return ServerConfig{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	c := ServerConfig{
		serverAddr: strings.TrimSpace(cfg.ServerAddr),
	}

	return c, nil
}

// ServerAddr returns the network address the server will listen on.
func (c ServerConfig) ServerAddr() string {
	return c.serverAddr
}

// ToOptions converts the configuration to a list of server options.
func (c ServerConfig) ToOptions() []ServerOption {
	var options []ServerOption
	if c.serverAddr != "" {
		options = append(options, WithServerAddr(c.serverAddr))
	}
	return options
}
