package mpserver

import (
	"os"
	"time"

	"arcadium.dev/core/build"
)

func (s MultiprotocolServer) GetLogLevel() string               { return s.loglevel }
func (s MultiprotocolServer) GetShutdownTimeout() time.Duration { return s.shutdownTimeout }
func (s MultiprotocolServer) GetServers() []ProtocolServer      { return s.servers }
func (s MultiprotocolServer) GetBuildInfo() build.Information   { return s.info }
func (s MultiprotocolServer) GetInterrupt() chan os.Signal      { return s.interrupt }
