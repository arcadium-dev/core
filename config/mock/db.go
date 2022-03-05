package mock

// DB is a mock of the db config.
type DB struct {
	DSN_, DriverName_ string
}

// DriverName returns the mocked db driver name.
func (d DB) DriverName() string {
	return d.DriverName_
}

// DSN returns the mocked db DSN.
func (d DB) DSN() string {
	return d.DSN_
}
