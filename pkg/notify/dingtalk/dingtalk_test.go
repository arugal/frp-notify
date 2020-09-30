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
	"github/arugal/frp-notify/pkg/config"
	"reflect"
	"testing"
)

func Test_parseAndVerifyConfig(t *testing.T) {
	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantConfig config.DingTalkConfig
		wantErr    bool
	}{
		{
			"normal",
			args{
				cfg: map[string]interface{}{
					"token":     "token",
					"secret":    "secret",
					"is_at_all": true,
				},
			},
			config.DingTalkConfig{
				Token:   "token",
				Secret:  "secret",
				IsAtAll: true,
			},
			false,
		},
		{
			"miss token",
			args{
				cfg: map[string]interface{}{},
			},
			config.DingTalkConfig{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := parseAndVerifyConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAndVerifyConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("parseAndVerifyConfig() gotConfig = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
