package server

import (
	"testing"

	"arcadium.dev/core/assert"
	"golang.org/x/crypto/ssh"
)

func TestWithAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		arg  string
		want func(t *testing.T, server *Server)
	}{
		{
			name: "empty address",
			want: func(t *testing.T, server *Server) {
				assert.Equal(t, server.addr, defaultAddr)
			},
		},
		{
			name: "valid address",
			arg:  ":2222",
			want: func(t *testing.T, server *Server) {
				assert.Equal(t, server.addr, ":2222")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := &Server{addr: defaultAddr}
			WithAddr(test.arg).apply(server)
			test.want(t, server)
		})
	}
}

func TestWithPublicKeyAuthn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		precond     PublicKeyCallback
		arg         PublicKeyCallback
		want        func(t *testing.T, server *Server)
		wantRecover func(t *testing.T, s string)
	}{
		{
			name: "nil callback",
			want: func(t *testing.T, server *Server) {
				assert.Nil(t, server.config.PublicKeyCallback)
			},
		},
		{
			name: "already set",
			precond: PublicKeyCallback(func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
				return nil, nil
			}),
			wantRecover: func(t *testing.T, s string) {
				assert.Equal(t, s, "public key callback already defined in ssh server config")
			},
		},
		{
			name: "valid callback",
			arg: PublicKeyCallback(func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
				return nil, nil
			}),
			want: func(t *testing.T, server *Server) {
				assert.NotNil(t, server.config.PublicKeyCallback)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if s, ok := recover().(string); ok {
					test.wantRecover(t, s)
				}
			}()

			server := &Server{config: ssh.ServerConfig{PublicKeyCallback: test.precond}}
			WithPublicKeyAuthn(test.arg).apply(server)
			test.want(t, server)
		})
	}
}

func TestWithPasswordAuthn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		precond     PasswordCallback
		arg         PasswordCallback
		want        func(t *testing.T, server *Server)
		wantRecover func(t *testing.T, s string)
	}{
		{
			name: "nil callback",
			want: func(t *testing.T, server *Server) {
				assert.Nil(t, server.config.PasswordCallback)
			},
		},
		{
			name: "already set",
			precond: PasswordCallback(func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
				return nil, nil
			}),
			wantRecover: func(t *testing.T, s string) {
				assert.Equal(t, s, "password callback already defined in ssh server config")
			},
		},
		{
			name: "valid callback",
			arg: PasswordCallback(func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
				return nil, nil
			}),
			want: func(t *testing.T, server *Server) {
				assert.NotNil(t, server.config.PasswordCallback)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if s, ok := recover().(string); ok {
					test.wantRecover(t, s)
				}
			}()

			server := &Server{config: ssh.ServerConfig{PasswordCallback: test.precond}}
			WithPasswordAuthn(test.arg).apply(server)
			test.want(t, server)
		})
	}
}
