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

package log

import (
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.Log

	notify.RegisterNotify("log", logNotifyBuilder)
}

type logNotify struct {
}

func (l *logNotify) SendMessage(title string, message string) {
	log.Infof("title: %s, message: %s \n", title, message)
}

func logNotifyBuilder(config map[string]interface{}) (notify.Notify, error) {
	return &logNotify{}, nil
}
