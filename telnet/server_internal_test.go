package telnet

import (
	"github.com/globalcyberalliance/telnet-go"
	"github.com/rs/zerolog"
)

func (s *Server) Addr() string                   { return s.addr }
func (s *Server) Handle(session *telnet.Session) { s.handle(session) }
func (s *Server) Logger() *zerolog.Logger        { return s.logger }
