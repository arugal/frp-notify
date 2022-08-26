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

package handler

import (
	"fmt"
	"os"
	"strconv"
	"time"

	telegramLib "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github/arugal/frp-notify/pkg/controller"
	"github/arugal/frp-notify/pkg/ip"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
	"github/arugal/frp-notify/pkg/types"

	"github.com/sirupsen/logrus"
)

const (
	defaultWindowInterval = time.Hour
	maxRequestSize        = 3000

	localIP = "127.0.0.1"
)

var log *logrus.Logger

func init() {
	log = logger.Log
}

type notifyHandler struct {
	windowInterval time.Duration
	requestChan    chan *types.Request

	addressService ip.AddressService
}

func NewNotifyHandler(ops ...NotifyHandlerOption) controller.HandlerChain {
	handler := &notifyHandler{
		windowInterval: defaultWindowInterval,
		requestChan:    make(chan *types.Request, maxRequestSize),
	}

	for _, op := range ops {
		op(handler)
	}

	return handler
}

type NotifyHandlerOption func(m *notifyHandler)

func WithAddressService(enable bool, newServiceFunc func() ip.AddressService) NotifyHandlerOption {
	return func(m *notifyHandler) {
		if enable {
			m.addressService = newServiceFunc()
		}
	}
}

func (n *notifyHandler) Op(op string) bool {
	switch op {
	case types.OpLogin, types.OpNewUserConn, types.OpNewWorkConn, types.OpNewProxy, types.OpPing, types.OpCloseProxy, types.OpExit:
		return true
	default:
		return false
	}
}

func (n *notifyHandler) Do(req *types.Request) (bool, *types.Response) {
	select {
	case n.requestChan <- req:
	default:
		log.Warnf("reach max send buffer")
	}
	return true, nil
}

func (n *notifyHandler) Open() error {
	go n.doNotify()
	return nil
}

func (n *notifyHandler) Close() error {
	close(n.requestChan)
	return nil
}

func (n notifyHandler) doNotify() {
	// TODO code refactor
	windowTicker := time.NewTicker(n.windowInterval)
	userConnCache := make(map[string]map[string]bool)
	for {
		select {
		case request := <-n.requestChan:
			log.Infof("receive new request op: %s, content: %s \n", request.Op, request.Content)

			var title = request.Op
			var message string
			var skipNotify = false
			switch request.Op {
			case types.OpLogin:
				login := request.Body.(*types.Login)
				title = fmt.Sprintf("%s online", login.Metas["hostname"])
				message = fmt.Sprintf("_Version: %v_\n_Time: %v_\n_Client IP: %v_",
					login.Version, convertTimestampToDatetime(login.Timestamp), login.RemoteIP)
			case types.OpExit:
				exitMsg := request.Body.(*types.Exit)
				title = fmt.Sprintf("%s offline", exitMsg.Metas["hostname"])
				message = fmt.Sprintf("_Version: %v_\n_Time: %v_\n_Client IP: %v_",
					exitMsg.Version, convertTimestampToDatetime(exitMsg.Timestamp), exitMsg.RemoteIP)
			case types.OpNewProxy:
				proxy := request.Body.(*types.Proxy)
				message = fmt.Sprintf("ProxyName: %v, ProxyType: %v, RemotePort: %v",
					telegramLib.EscapeText(telegramLib.ModeMarkdown, proxy.ProxyName), proxy.ProxyType, proxy.RemotePort)
			case types.OpNewWorkConn:
				workConn := request.Body.(*types.WorkConn)
				message = fmt.Sprintf("RunID: %v", workConn.RunID)
			case types.OpNewUserConn:
				userConn := request.Body.(*types.UserConn)
				// ip cache
				proxyName := fmt.Sprint(userConn.ProxyName)
				ipCache, ok := userConnCache[proxyName]
				if !ok {
					ipCache = make(map[string]bool)
					userConnCache[proxyName] = ipCache
				}
				if _, ok := ipCache[userConn.RemoteIP]; ok {
					skipNotify = true
					break
				}
				ipCache[userConn.RemoteIP] = false
				var ipCName string
				if n.addressService != nil && userConn.RemoteIP != localIP {
					ipCName = n.addressService.Query(userConn.RemoteIP)
				}
				if ipCName != "" {
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s, CName: %s", telegramLib.EscapeText(telegramLib.ModeMarkdown, proxyName), userConn.ProxyType, userConn.RemoteIP, ipCName)
				} else {
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s", telegramLib.EscapeText(telegramLib.ModeMarkdown, proxyName), userConn.ProxyType, userConn.RemoteIP)
				}

			case types.OpPing:
				ping := request.Body.(*types.Ping)
				message = fmt.Sprintf("ping %s %s", ping.PrivilegeKey, ping.User.User)

			case types.OpCloseProxy:
				proxy := request.Body.(*types.CloseProxy)
				message = fmt.Sprintf("Closed proxy %s", telegramLib.EscapeText(telegramLib.ModeMarkdown, proxy.ProxyName))

			}
			if !skipNotify {
				notify.SendMessage(title, message)
			}
		case <-windowTicker.C:
			userConnCache = make(map[string]map[string]bool)
			log.Debugln("clean user conn cache.")
		}
	}
}

func convertTimestampToDatetime(timestamp int64) string {
	loc := time.UTC
	if name, ok := os.LookupEnv("TZ"); ok {
		if name != "" {
			loc, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
		}
	}
	ti, _ := strconv.ParseInt(strconv.FormatInt(timestamp, 10), 10, 64)
	tm := time.Unix(ti, 0)
	return tm.In(loc).Format("Jan 2, 2006 at 3:04 PM")
}
