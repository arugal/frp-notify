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
	"github/arugal/frp-notify/pkg/controller"
	"github/arugal/frp-notify/pkg/ip"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
	"github/arugal/frp-notify/pkg/types"
	"time"

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
	case types.OpLogin, types.OpNewUserConn, types.OpNewWorkConn, types.OpNewProxy:
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
				message = fmt.Sprintf("Version: %v, HostName: %v, Os: %v, Arch: %v",
					login.Version, login.Hostname, login.OS, login.Arch)
			case types.OpNewProxy:
				proxy := request.Body.(*types.Proxy)
				message = fmt.Sprintf("ProxyName: %v, ProxyType: %v, RemotePort: %v",
					proxy.ProxyName, proxy.ProxyType, proxy.RemotePort)
			case types.OpCloseProxy:
				closeProxy := request.Body.(*types.CloseProxy)
				message = fmt.Sprintf("ProxyName: %v", closeProxy.ProxyName)
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
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s, CName: %s", proxyName, userConn.ProxyType, userConn.RemoteIP, ipCName)
				} else {
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s", proxyName, userConn.ProxyType, userConn.RemoteIP)
				}
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
