//  Copyright 2026 arcadium.dev <info@arcadium.dev>
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package telnet

import (
	"runtime/debug"

	"github.com/rs/zerolog"
)

type (
	RecoveryMiddleware struct {
		Logger *zerolog.Logger
	}
)

func (r RecoveryMiddleware) Recover(next Handler) Handler {
	return HandlerFunc(func(session *Session) {
		defer func() {
			if recover() == nil {
				return
			}
			r.Logger.Error().Msg("recovering from a panic")
			r.Logger.Error().Msgf("stacktrace: %s", debug.Stack())
		}()
		next.ServeTELNET(session)
	})
}

type (
	SessionMiddleware struct {
		Logger *zerolog.Logger
	}
)

func (s SessionMiddleware) Session(next Handler) Handler {
	return HandlerFunc(func(session *Session) {
		s.Logger.Debug().Str("remote addr", session.RemoteAddr().String()).Msg("session start")
		defer s.Logger.Debug().Str("remote addr", session.RemoteAddr().String()).Msg("session complete")
		next.ServeTELNET(session)
	})
}
