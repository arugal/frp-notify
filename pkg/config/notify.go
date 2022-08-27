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

package config

// NotifyConfig notify configuration
type NotifyConfig struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// GotifyConfig gotify configuration
type GotifyConfig struct {
	ServerProto string `json:"server_proto"`
	ServerAddr  string `json:"server_addr"`
	AppToken    string `json:"app_token"`
}

// DingTalkConfig dingTalk configuration
type DingTalkConfig struct {
	Token   string `json:"token"`
	Secret  string `json:"secret"`
	IsAtAll bool   `json:"is_at_all"`
}

// WxWorkConfig wxwork configuration
type WxWorkConfig struct {
	CorpID       string   `json:"corp_id"`
	CorpSecret   string   `json:"corp_secret"`
	AgentID      int64    `json:"agent_id"`
	ToUser       []string `json:"to_user"`
	ToParty      []string `json:"to_party"`
	ToTag        []string `json:"to_tag"`
	FilterRegExp string   `json:"filter_regexp"`
	AdminURL     string   `json:"admin_url"`
}

type LarkConfig struct {
	WebhookURL string   `json:"webhook_url"`
	Secret     string   `json:"secret"`
	AtUsers    []string `json:"at_users"`
}
