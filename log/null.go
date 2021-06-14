// Copyright 2021 arcadium.dev <info@arcadium.dev>
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

func NewNullLogger() NullLogger { return NullLogger{} }

type NullLogger struct{}

func (l NullLogger) WithField(string, interface{}) Logger     { return l }
func (l NullLogger) WithFields(map[string]interface{}) Logger { return l }
func (l NullLogger) WithError(err error) Logger               { return l }

func (l NullLogger) Debug(...interface{})          {}
func (l NullLogger) Debugln(...interface{})        {}
func (l NullLogger) Debugf(string, ...interface{}) {}

func (l NullLogger) Info(...interface{})          {}
func (l NullLogger) Infoln(...interface{})        {}
func (l NullLogger) Infof(string, ...interface{}) {}

func (l NullLogger) Print(...interface{})          {}
func (l NullLogger) Println(...interface{})        {}
func (l NullLogger) Printf(string, ...interface{}) {}

func (l NullLogger) Warning(...interface{})          {}
func (l NullLogger) Warningln(...interface{})        {}
func (l NullLogger) Warningf(string, ...interface{}) {}

func (l NullLogger) Error(...interface{})          {}
func (l NullLogger) Errorln(...interface{})        {}
func (l NullLogger) Errorf(string, ...interface{}) {}

func (l NullLogger) Fatal(...interface{})          {}
func (l NullLogger) Fatalln(...interface{})        {}
func (l NullLogger) Fatalf(string, ...interface{}) {}

func (l NullLogger) V(level int) bool { return false }
