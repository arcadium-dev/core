package server

import (
	"testing"

	"arcadium.dev/core/assert"
)

func TestNewSessions(t *testing.T) {
	assert.NotNil(t, NewSessions())
}

func TestSessionsAdd(t *testing.T) {
	t.Run("close failure", func(t *testing.T) {
		sessions := &Sessions{closing: true}
		assert.Error(t, sessions.Add(nil), `cannot add session, closing`)
	})

	t.Run("success", func(t *testing.T) {
		sessions := NewSessions()
		session := &Session{}
		assert.Nil(t, sessions.Add(session))
		assert.Equal(t, len(sessions.sessions), 1)
		_, ok := sessions.sessions[session]
		assert.True(t, ok)
	})
}

func TestSessionsRemove(t *testing.T) {
	t.Run("close failure", func(t *testing.T) {
		sessions := &Sessions{closing: true}
		assert.Error(t, sessions.Remove(nil), `cannot remove session, closing`)
	})

	t.Run("success", func(t *testing.T) {
		sessions := NewSessions()
		session := &Session{}
		assert.Nil(t, sessions.Add(session))
		assert.Nil(t, sessions.Remove(session))
		assert.Equal(t, len(sessions.sessions), 0)
		_, ok := sessions.sessions[session]
		assert.False(t, ok)
	})

}

func TestSessionsClose(t *testing.T) {
	sessions := NewSessions()
	var s []*Session
	for i := 0; i < 10; i++ {
		s = append(s, &Session{})
	}
	for _, session := range s {
		assert.Nil(t, sessions.Add(session))
	}
	sessions.Close()
	sessions.Wait()
}
