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

type FRPNotifyConfig struct {
	BindAddress    string         `json:"-"`
	WindowInterval int64          `json:"-"`
	Blacklist      []string       `json:"blacklist"`
	Whitelist      []string       `json:"whitelist"` // If a handler is configured, only the IP within the handler can be accessed
	NotifyPlugins  []NotifyConfig `json:"notify_plugins"`
}
