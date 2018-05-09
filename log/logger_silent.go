// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright 2018 Banco Bilbao Vizcaya Argentaria, S.A.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"log"
)

type silentLogger struct {
	log.Logger
}

func newSilent() *silentLogger {
	return &silentLogger{}
}

// A impl 'l Nologger' qed/log.Logger
func (l silentLogger) Error(v ...interface{})                 { return }
func (l silentLogger) Warn(v ...interface{})                  { return }
func (l silentLogger) Info(v ...interface{})                  { return }
func (l silentLogger) Debug(v ...interface{})                 { return }
func (l silentLogger) Errorf(format string, v ...interface{}) { return }
func (l silentLogger) Warnf(format string, v ...interface{})  { return }
func (l silentLogger) Infof(format string, v ...interface{})  { return }
func (l silentLogger) Debugf(format string, v ...interface{}) { return }
