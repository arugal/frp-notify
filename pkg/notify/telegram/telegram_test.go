package telegram

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
		wantConfig config.TelegramConfig
		wantErr    bool
	}{
		{
			"normal",
			args{
				cfg: map[string]interface{}{
					"token":   "token",
					"groupId": 10,
				},
			},
			config.TelegramConfig{
				Token:   "token",
				GroupId: 10,
			},
			false,
		},
		{
			"miss token",
			args{
				cfg: map[string]interface{}{},
			},
			config.TelegramConfig{},
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

func Test_SendMessage(t *testing.T) {
	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		message string
	}{
		{
			"normal",
			args{
				cfg: map[string]interface{}{
					"token":   "54...:AAE...",
					"groupId": -100,
				},
			},
			"This is a test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notify, err := telegramNotifyBuilder(tt.args.cfg)
			if err != nil {
				return
			}
			notify.SendMessage("Test", "detail")

		})
	}
}
