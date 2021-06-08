// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server // import "arcadium.dev/core/server"

//go:generate mockgen -package mockserver -destination ./mock/server.go . Server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// Server abstracts a server.
type Server interface {
	// Serve is the entry point for a service and will be run it a goroutine.
	// It is passed a channel to communicate the result.
	Serve(result chan<- error)

	// Stop stops the server.
	Stop()
}

// Serve starts the server and will catch os signals. If an os signal is
// caught, the server will be cancelled via it's the context.
func Serve(s Server) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	result := make(chan error, 1)
	return serve(sig, result, s)
}

func serve(sig chan os.Signal, result chan error, server Server) error {
	go server.Serve(result)

	var err error
	select {
	case err = <-result:
		if err != nil {
			log.Printf("\n\nError: %s\n\n", err.Error())
			if err, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
				for _, f := range err.StackTrace() {
					log.Printf("%+s:%d\n", f, f)
				}
			}
		}
	case signal := <-sig:
		log.Printf("\n\nsignal received: %s\n\n", signal)
	}

	server.Stop()
	time.Sleep(1 * time.Second) // Give the logs some time to flush.

	return err
}
