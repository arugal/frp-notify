// Copyright 2020 arugal, zhangwei24@apache.org
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

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	if Log == nil {
		Log = logrus.New()
	}
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
}

func SetLogLevel(level string) {
	switch level {
	case "panic":
		Log.SetLevel(logrus.PanicLevel)
	case "fatal":
		Log.SetLevel(logrus.FatalLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "trace":
		Log.SetLevel(logrus.DebugLevel)
	}
}
