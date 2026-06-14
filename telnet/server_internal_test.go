package telnet

import (
	"log/slog"

	"arcadium.dev/telnet"
)

func (s Server) Addr() string                   { return s.addr }
func (s Server) Logger() *slog.Logger           { return s.logger }
func (s Server) Handle(session *telnet.Session) { s.handle(session) }
