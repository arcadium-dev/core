package mpserver

import (
	"os"
	"time"

	"arcadium.dev/core/build"
)

func (s MultiprotocolServer) GetLogLevel() string               { return s.loglevel }
func (s MultiprotocolServer) GetShutdownTimeout() time.Duration { return s.shutdownTimeout }
func (s MultiprotocolServer) GetBuildInfo() build.Information   { return s.info }
func (s MultiprotocolServer) GetInterrupt() chan os.Signal      { return s.interrupt }

type Servers = servers

func (s MultiprotocolServer) Servers() *Servers {
	return s.servers
}

func (s *Servers) Len() int                 { return s.len() }
func (s *Servers) Get(i int) ProtocolServer { return s.get(i) }
