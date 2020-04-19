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

package notify

import (
	"errors"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"

	"github.com/sirupsen/logrus"
)

var (
	notifys = make(map[string]*holder)

	log *logrus.Logger

	errNoSuchNotify = errors.New("no such notify")
)

func init() {
	log = logger.Log
}

type holder struct {
	notify        Notify
	inited        bool
	notifyBuilder PluginBuilder
}

func InitNotify(cfg config.NotifyConfig) (err error) {
	holder, ok := notifys[cfg.Name]
	if ok {
		holder.notify, err = holder.notifyBuilder(cfg.Config)
		if err == nil {
			holder.inited = true
		}
		return err
	}
	return errNoSuchNotify
}

func SendMessage(title string, message string) {
	for name, holder := range notifys {
		if holder.inited {
			log.Debugf("send %s:%s to %s.", title, message, name)
			holder.notify.SendMessage(title, message)
		}
	}
}

func RegisterNotify(name string, builder PluginBuilder) {
	notifys[name] = &holder{
		notifyBuilder: builder,
	}
}
