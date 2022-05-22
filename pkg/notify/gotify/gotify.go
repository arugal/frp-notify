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

package gotify

import (
	"encoding/json"
	"fmt"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.Log

	notify.RegisterNotify("gotify", gotifyNotifyBuilder)
}

type gotifyNotify struct {
	cfg    config.GotifyConfig
	client *resty.Client
}

func (g *gotifyNotify) SendMessage(title string, message string) {
	format := make(map[string]string)
	format["title"] = title
	format["message"] = message

	_, err := g.client.R().
		SetFormData(format).
		Post(fmt.Sprintf("%s://%s/message?token=%s", g.cfg.ServerProto, g.cfg.ServerAddr, g.cfg.AppToken))
	if err != nil {
		log.Errorf("send message to gotify error, err: %s serverProto: %s, serverAddr: %s, token: %s", err,
			g.cfg.ServerProto, g.cfg.ServerAddr, g.cfg.AppToken)
	}
}

func parseAndVerifyConfig(cfg map[string]interface{}) (config config.GotifyConfig, err error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return config, err
	}
	log.Debugf("gotify config: %v", string(data))

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}
	if config.ServerAddr == "" {
		return config, fmt.Errorf("miss server_addr")
	}
	if config.AppToken == "" {
		return config, fmt.Errorf("miss app_token")
	}
	return
}

func gotifyNotifyBuilder(cfg map[string]interface{}) (notify.Notify, error) {
	gotifyConfig, err := parseAndVerifyConfig(cfg)
	if err != nil {
		return nil, err
	}
	client := resty.New().SetTimeout(time.Second * 3)
	return &gotifyNotify{
		cfg:    gotifyConfig,
		client: client,
	}, nil
}
