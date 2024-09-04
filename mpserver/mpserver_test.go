package mpserver_test

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/build"
	"arcadium.dev/core/log"
	"arcadium.dev/core/mpserver"
	"arcadium.dev/core/require"
)

const (
	version = "v0.0.1"
	branch  = "dev-branch"
	commit  = "commit-sha-goes-here"
	date    = "today"
)

func TestMPServer_New(t *testing.T) {
	tests := []struct {
		name   string
		opts   []mpserver.Option
		verify func(t *testing.T, s *mpserver.MultiprotocolServer, err error)
	}{
		{
			name: "log creation failure",
			opts: []mpserver.Option{
				mpserver.WithLogLevel("invalid log level"),
			},
			verify: func(t *testing.T, s *mpserver.MultiprotocolServer, err error) {
				assert.Nil(t, s)
				assert.Error(t, err, "failed to create logger: invalid level: 6")
			},
		},
		{
			name: "success",
			opts: []mpserver.Option{
				mpserver.WithLogLevel("debug"),
				mpserver.WithShutdownTimeout(600 * time.Second),
			},
			verify: func(t *testing.T, s *mpserver.MultiprotocolServer, err error) {
				assert.Nil(t, err)
				require.NotNil(t, s)
				assert.Equal(t, s.GetLogLevel(), "debug")
				assert.Equal(t, s.GetShutdownTimeout(), 600*time.Second)
				assert.Equal(t, s.GetBuildInfo(), build.Information{
					Name:    "mpserver.test",
					Version: version,
					Branch:  branch,
					Commit:  commit,
					Date:    date,
					Go:      runtime.Version(),
				})
				assert.NotNil(t, s.GetInterrupt)
				assert.NotNil(t, s.Ctx())
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			s, err := mpserver.New(version, branch, commit, date, test.opts...)
			test.verify(t, s, err)
		})
	}
}

func TestMPServer_Register(t *testing.T) {
	var (
		version = "version"
		branch  = "branch"
		commit  = "commit"
		date    = "date"
	)

	ctx, b := log.SetupTestLogging(t)
	s, err := mpserver.New(version, branch, commit, date)
	assert.Nil(t, err)

	server := mockProtocolServer{}

	s.Register(ctx, server)

	require.Equal(t, b.Len(), 1)
	assert.Equal(t, b.Index(0), `{"severity":"info","message":"protocol server registered: mockProtocolServer"}`+"\n")
}

func TestMPServer_Serve(t *testing.T) {
	tests := []struct {
		name    string
		opts    []mpserver.Option
		servers []mpserver.ProtocolServer
		verify  func(t *testing.T, err error)
	}{
		{
			name: "nothing to serve failure",
			opts: []mpserver.Option{
				mpserver.WithLogLevel("error"),
				mpserver.WithShutdownTimeout(300 * time.Second),
			},
			verify: func(t *testing.T, err error) {
				assert.Error(t, err, "exiting, nothing to server")
			},
		},
		{
			name: "protocol server fails to start",
			opts: []mpserver.Option{
				mpserver.WithLogLevel("error"),
				mpserver.WithShutdownTimeout(1 * time.Second),
			},
			servers: []mpserver.ProtocolServer{
				mockProtocolServer{err: errors.New("failed to start, 1")},
				mockProtocolServer{err: errors.New("failed to start, 2")},
				mockProtocolServer{err: errors.New("failed to start, 3")},
			},
			verify: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "failed to start, ")
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			s, err := mpserver.New(version, branch, commit, date, test.opts...)
			require.Nil(t, err)
			s.Register(s.Ctx(), test.servers...)
			err = s.Serve()
			test.verify(t, err)
		})
	}
}

func TestMPServer_Shutdown(t *testing.T) {
	s, err := mpserver.New(version, branch, commit, date)
	require.Nil(t, err)
	s.Register(s.Ctx(), []mpserver.ProtocolServer{
		mockProtocolServer{},
		mockProtocolServer{},
		mockProtocolServer{},
		mockProtocolServer{},
	}...)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = s.Serve()
		assert.Nil(t, err)
		wg.Done()
	}()
	s.Shutdown()
	wg.Wait()
}

type (
	mockProtocolServer struct {
		err error
	}
)

func (m mockProtocolServer) Serve(context.Context) error {
	return m.err
}

func (m mockProtocolServer) Shutdown(context.Context) {}

func (m mockProtocolServer) Name() string { return "mockProtocolServer" }
