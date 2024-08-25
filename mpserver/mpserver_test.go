package mpserver_test

import (
	"context"
	"testing"

	"arcadium.dev/core/build"
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
			verify: func(t *testing.T, s *mpserver.MultiprotocolServer, err error) {
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

func TestMPServer_Serve(t *testing.T) {
	tests := []struct {
		name   string
		opts   []mpserver.Option
		verify func(t *testing.T, err error)
	}{
		{
			name: "nothing to server failure",
			verify: func(t *testing.T, err error) {
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			s, err := mpserver.New(version, branch, commit, date, test.opts...)
			require.Nil(t, err)

			err = s.Serve()
			test.verify(t, err)
		})
	}
}

type (
	mockProtocolServer struct {
		err error
	}
)

func (m mockProtocolServer) Serve(context.Context, build.Information) error {
	return m.err
}

func (m mockProtocolServer) Shutdown(context.Context) {}
