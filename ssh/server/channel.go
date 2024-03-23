// Copyright 2023-2024 arcadium.dev <info@arcadium.dev>
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

package server // import "arcadium.dev/core/ssh/server"

import (
	"golang.org/x/crypto/ssh"
)

// ChannelType as defined in https://www.rfc-editor.org/rfc/rfc4250.html#section-4.9.1
type ChannelType string

const (
	ChannelTypeSession        = ChannelType("session")
	ChannelTypeX11            = ChannelType("x11")
	ChannelTypeForwardedTCPIP = ChannelType("forwarded-tcpip")
	ChannelTypeDirectTCPIP    = ChannelType("direct-tcpip")
)

// ChannelHandler defines the behavior of a service that handles a specific channel type.
type ChannelHandler interface {
	Type() ChannelType
	Handle(*ssh.ServerConn, ssh.NewChannel)
	Shutdown() error
	Close() error
}
