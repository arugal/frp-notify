package telegram

import (
	"encoding/json"
	"fmt"
	telegramLib "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
)

var (
	log *logrus.Logger
)

type telegramNotify struct {
	cfg config.TelegramConfig
	api *telegramLib.BotAPI
}

func init() {
	log = logger.Log
	notify.RegisterNotify("telegram", telegramNotifyBuilder)
}

func parseAndVerifyConfig(cfg map[string]interface{}) (config config.TelegramConfig, err error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return config, err
	}
	log.Debugf("telegram config: %v", string(data))

	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	if config.Token == "" {
		return config, fmt.Errorf("miss token")
	}
	if config.GroupId == 0 {
		return config, fmt.Errorf("miss group id")
	}
	return config, nil
}

func telegramNotifyBuilder(cfg map[string]interface{}) (notify.Notify, error) {
	telegramConfig, err := parseAndVerifyConfig(cfg)
	if err != nil {
		return nil, err
	}

	client, err := telegramLib.NewBotAPI(telegramConfig.Token)
	if err != nil {
		log.Panic(err)
	}

	client.Debug = true

	log.Printf("Authorized Telegram on account %s", client.Self.UserName)

	return &telegramNotify{
		cfg: telegramConfig,
		api: client,
	}, nil
}

func (t *telegramNotify) SendMessage(title string, message string) {
	msg := telegramLib.NewMessage(t.cfg.GroupId, fmt.Sprintf("*FRP Server* said: %s \n_%s_", title, message))
	msg.ParseMode = telegramLib.ModeMarkdown
	_, err := t.api.Send(msg)
	if err != nil {
		log.Errorf("send message to telegram error, err: %s", err)
	}
}
