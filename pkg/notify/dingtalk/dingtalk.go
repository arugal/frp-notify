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

package dingtalk

import (
	"encoding/json"
	"fmt"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"

	robot "github.com/JetBlink/dingtalk-notify-go-sdk"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	log = logger.Log

	notify.RegisterNotify("dingTalk", dingTalkNotifyBuilder)
}

type dingTalkNotify struct {
	cfg      config.DingTalkConfig
	dingTalk *robot.Robot
}

func (d *dingTalkNotify) SendMessage(title string, message string) {
	err := d.dingTalk.SendMarkdownMessage("dingTalk", fmt.Sprintf("### %s \n\n %s", title, message), []string{}, d.cfg.IsAtAll)
	if err != nil {
		log.Errorf("send message to dingTalk error, err: %s", err)
	}
}

func parseAndVerifyConfig(cfg map[string]interface{}) (config config.DingTalkConfig, err error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return config, err
	}
	log.Debugf("dingTalk config: %v", string(data))

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	if config.Token == "" {
		return config, fmt.Errorf("miss token")
	}
	return
}

func dingTalkNotifyBuilder(cfg map[string]interface{}) (notify.Notify, error) {
	dingTalkConfig, err := parseAndVerifyConfig(cfg)
	if err != nil {
		return nil, err
	}
	client := robot.NewRobot(dingTalkConfig.Token, dingTalkConfig.Secret)
	return &dingTalkNotify{
		cfg:      dingTalkConfig,
		dingTalk: client,
	}, nil
}
