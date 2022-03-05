package mock

// Logger is a mock of the Logger config.
type Logger struct {
	Level_, File_, Format_ string
}

// Level returns the mocked logging level.
func (m Logger) Level() string {
	return m.Level_
}

// File returns the mocked logging output file.
func (m Logger) File() string {
	return m.File_
}

// Format returns the mocked logging format.
func (m Logger) Format() string {
	return m.Format_
}
