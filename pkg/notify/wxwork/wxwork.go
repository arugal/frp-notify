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

package wxwork

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xen0n/go-workwx"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
	"regexp"
)

var (
	log *logrus.Logger
)

type wxworkNotify struct {
	cfg config.WxWorkConfig
	app *workwx.WorkwxApp
}

func init() {
	log = logger.Log

	notify.RegisterNotify("wxwork", wxworkNotifyBuilder)
}

func parseAndVerifyConfig(cfg map[string]interface{}) (config config.WxWorkConfig, err error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return config, err
	}
	log.Debugf("wxwork config: %v", string(data))

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	if config.CorpId == "" {
		return config, fmt.Errorf("miss corp_id")
	}
	if config.CorpSecret == "" {
		return config, fmt.Errorf("miss corp_secret")
	}
	if config.AgentId == 0 {
		return config, fmt.Errorf("miss agent_id")
	}

	if len(config.ToUser) == 0 && len(config.ToParty) == 0 && len(config.ToTag) == 0 {
		config.ToUser = []string{"@all"}
	}

	if config.FilterRegExp == "" {
		config.FilterRegExp = ".*"
	}

	return
}

func wxworkNotifyBuilder(cfg map[string]interface{}) (notify.Notify, error) {
	wxworkdConfig, err := parseAndVerifyConfig(cfg)
	if err != nil {
		return nil, err
	}

	client := workwx.New(wxworkdConfig.CorpId)

	app := client.WithApp(wxworkdConfig.CorpSecret, wxworkdConfig.AgentId)
	app.SpawnJSAPITicketRefresher()

	return &wxworkNotify{
		cfg: wxworkdConfig,
		app: app,
	}, nil
}

func (wx *wxworkNotify) SendMessage(title string, message string) {
	toWho := workwx.Recipient{
		UserIDs:  wx.cfg.ToUser,
		PartyIDs: wx.cfg.ToParty,
		TagIDs:   wx.cfg.ToTag,
	}

	filterRegexp, err := regexp.Compile(wx.cfg.FilterRegExp)

	if err != nil {
		log.Errorf("wxwork error msg filter: %s, err: %s", wx.cfg.FilterRegExp, err)
	}

	if filterRegexp.MatchString(message) {
		err = wx.app.SendTextCardMessage(&toWho, title, message, wx.cfg.AdminUrl, "后台详情", false)

		if err != nil {
			log.Errorf("send message to wxwork error, err: %s", err)
		}
	}
}
