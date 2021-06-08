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

package log // import "arcadium.dev/core/log"

//go:generate mockgen -package mocklog -destination ./mock/logger.go . Logger

// Logger implements grpc's LoggerV2 interface while supporting structured logging
// by implementing an interface similar to logrus' FieldLogger or apex's Interface.
//
// See:
//   - https://github.com/grpc/grpc-go/blob/v1.27.1/grpclog/loggerv2.go
//   - https://github.com/sirupsen/logrus/blob/v1.4.2/logrus.go
//   - https://github.com/apex/log/blob/v1.1.2/interface.go
type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger

	Debug(args ...interface{})
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})

	Print(args ...interface{})
	Println(args ...interface{})
	Printf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})

	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})

	V(l int) bool
}
