package lark

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
)

var (
	log *logrus.Logger
)

type larkNotify struct {
	cfg    *config.LarkConfig
	client *http.Client
}

type larkMessage struct {
	MsgType   string                 `json:"msg_type"`
	Content   map[string]interface{} `json:"content"`
	Timestamp string                 `json:"timestamp,omitempty"`
	Sign      string                 `json:"sign,omitempty"`
}

func (lark *larkNotify) buildMessage(title string, message string) *larkMessage {
	content := make([]map[string]interface{}, 0, len(lark.cfg.AtUsers)+1)
	content = append(content, map[string]interface{}{
		"tag":  "text",
		"text": message,
	})
	for _, atUser := range lark.cfg.AtUsers {
		content = append(content, map[string]interface{}{
			"tag":     "at",
			"user_id": atUser,
		})
	}
	msg := &larkMessage{
		MsgType: "post",
		Content: map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": title,
					"content": [][]map[string]interface{}{
						content,
					},
				},
			},
		},
	}

	if lark.cfg.Secret != "" {
		t := time.Now().Unix()
		msg.Timestamp = strconv.FormatInt(t, 10)
		msg.Sign = lark.sign(t)
	}
	return msg
}

func (lark *larkNotify) sign(t int64) string {
	s := fmt.Sprintf("%d\n%s", t, lark.cfg.Secret)
	h := hmac.New(sha256.New, []byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type sendMessageResponse struct {
	StatusCode *int `json:"StatusCode"`
}

func (lark *larkNotify) SendMessage(title string, message string) {
	msg := lark.buildMessage(title, message)
	b, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("marshal lark message failed: %s", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, lark.cfg.WebhookURL, bytes.NewReader(b))
	if err != nil {
		log.Errorf("new lark message request failed: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	resp, err := lark.client.Do(req)
	if err != nil {
		log.Errorf("send lark message request failed: %s", err)
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read lark message response failed: %s", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("send lark message failed, status: %d, resp: %s", resp.StatusCode, string(b))
		return
	}
	var res *sendMessageResponse
	if err = json.Unmarshal(b, &res); err != nil {
		log.Errorf("unmarshal lark message response failed: %s", err)
		return
	}
	if res.StatusCode == nil || *res.StatusCode != 0 {
		log.Errorf("send lark message failed, status: %d, resp: %s", *res.StatusCode, string(b))
	}
}

func init() {
	log = logger.Log

	notify.RegisterNotify("lark", larkNotifyBuilder)
}

func parseAndVerifyConfig(cfg map[string]interface{}) (*config.LarkConfig, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	log.Debugf("lark config: %s", string(data))

	var conf *config.LarkConfig
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	if conf.WebhookURL == "" {
		return nil, fmt.Errorf("miss webhook_url")
	}
	return conf, nil
}

func larkNotifyBuilder(cfg map[string]interface{}) (notify.Notify, error) {
	larkConfig, err := parseAndVerifyConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &larkNotify{
		cfg: larkConfig,
		client: &http.Client{
			Timeout: 600 * time.Second,
		},
	}, nil
}
