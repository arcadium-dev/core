package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s Server) Addr() string                   { return s.addr }
func (s Server) TLSConfig() *tls.Config         { return s.server.TLSConfig }
func (s Server) CORSOptions() *cors.Options     { return s.corsOptions }
func (s Server) ReadTimeout() time.Duration     { return s.server.ReadTimeout }
func (s Server) WriteTimeout() time.Duration    { return s.server.WriteTimeout }
func (s Server) ShutdownTimeout() time.Duration { return s.shutdownTimeout }

func (s Server) Router() *mux.Router {
	return s.router
}
func (s Server) Server() *http.Server {
	return s.server
}

type Services = services

func (s Server) Services() *Services {
	return s.services
}

func (s *Services) Len() int          { return s.len() }
func (s *Services) Get(i int) Service { return s.get(i) }
