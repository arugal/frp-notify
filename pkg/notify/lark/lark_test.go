package lark

import (
	"reflect"
	"testing"

	"github/arugal/frp-notify/pkg/config"
)

func Test_parseAndVerifyConfig(t *testing.T) {
	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *config.LarkConfig
		wantErr    bool
	}{
		{
			"normal",
			args{
				cfg: map[string]interface{}{
					"webhook_url": "webhook_url",
					"secret":      "secret",
				},
			},
			&config.LarkConfig{
				WebhookURL: "webhook_url",
				Secret:     "secret",
			},
			false,
		},
		{
			"miss webhook_url",
			args{
				cfg: map[string]interface{}{},
			},
			nil,
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
