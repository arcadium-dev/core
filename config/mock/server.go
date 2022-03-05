package mock

// Server is a mock of the Server config.
type Server struct {
	Addr_, Cert_, Key_, CACert_ string
}

// Addr returns the mocked server address.
func (s Server) Addr() string {
	return s.Addr_
}

// Cert returns the mocked server certificate.
func (s Server) Cert() string {
	return s.Cert_
}

// Key returns the mocked server certificate private key.
func (s Server) Key() string {
	return s.Key_
}

// CACert returns the mocked server CA cert.
func (s Server) CACert() string {
	return s.CACert_
}
