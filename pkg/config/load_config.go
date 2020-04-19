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

import (
	"encoding/json"
	"github/arugal/frp-notify/pkg/logger"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	config         FRPNotifyConfig
	configLoadFunc LoadFunc = DefaultConfigLoad

	log *logrus.Logger
)

func init() {
	log = logger.Log
}

type LoadFunc func(path string) *FRPNotifyConfig

func DefaultConfigLoad(path string) *FRPNotifyConfig {
	log.Infof("load config from : %s \n", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("load config failed, ", err)
	}
	cfg := &FRPNotifyConfig{}
	// translate to lower case
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalln("json unmarshal config failed, ", err)
	}
	return cfg
}

// Load config file and parse
func Load(path string) *FRPNotifyConfig {
	if cfg := configLoadFunc(path); cfg != nil {
		config = *cfg
	}
	return &config
}
