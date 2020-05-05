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

package server

import (
	"encoding/json"
	"fmt"
	"github/arugal/frp-notify/pkg/ip"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/notify"
	"github/arugal/frp-notify/pkg/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	defaultServerAddr     = ":80"
	defaultWindowInterval = time.Hour
	maxRequestSize        = 3000
)

var log *logrus.Logger

var normalResponse = types.Response{
	Reject:       false,
	RejectReason: "",
	Unchange:     true,
	Content:      nil,
}

func init() {
	log = logger.Log
}

type ManagerServer struct {
	serverAddr     string
	windowInterval time.Duration
	requestChan    chan *types.Request
	addressService ip.AddressService
}

func NewManagerServer(opts ...ManagerServerOption) *ManagerServer {
	ms := &ManagerServer{
		serverAddr:     defaultServerAddr,
		windowInterval: defaultWindowInterval,
		requestChan:    make(chan *types.Request, maxRequestSize),
	}

	for _, o := range opts {
		o(ms)
	}
	return ms
}

type ManagerServerOption func(m *ManagerServer)

// WithWindowInterval unit: minute
func WithWindowInterval(windowInterval int64) ManagerServerOption {
	return func(m *ManagerServer) {
		m.windowInterval = time.Duration(windowInterval) * time.Minute
	}
}

func WithServerAddr(serverAddr string) ManagerServerOption {
	return func(m *ManagerServer) {
		m.serverAddr = serverAddr
	}
}

func WithIPAddressService(service ip.AddressService) ManagerServerOption {
	return func(m *ManagerServer) {
		m.addressService = service
	}
}

// Start start manager service
func (m *ManagerServer) Start() {
	go m.doNotify()
	m.httpServer()
}

// Close
func (m *ManagerServer) Close() {
	close(m.requestChan)
}

func (m *ManagerServer) httpServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/handler", func(ctx *gin.Context) {
		defer func() {
			ctx.JSON(http.StatusOK, normalResponse)
		}()
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Warnf("read request body error: %v \n", err)
			return
		}

		request := &types.Request{}
		err = json.Unmarshal(body, request)
		if err != nil {
			log.Warnf("unmarshal request body error: %v \n", err)
			return
		}

		if request.Version != types.APIVersion {
			log.Warnf("unsupported api version %s \n", request.Version)
			return
		}
		if request.Op == types.OpPing {
			log.Debug("ignore ping operation")
			return
		}
		select {
		case m.requestChan <- request:
		default:
			log.Warnf("reach max send buffer")
		}
	})

	err := r.Run(m.serverAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *ManagerServer) doNotify() {
	windowTicker := time.NewTicker(m.windowInterval)
	userConnCache := make(map[string]map[string]bool)
	for {
		select {
		case request := <-m.requestChan:
			log.Infof("receive new request op: %s, content: %s \n", request.Op, request.Content)

			content, ok := request.Content.(map[string]interface{})
			if !ok {
				log.Warnf("unsupported content: %s \n", request.Content)
				return
			}
			var title = request.Op
			var message string
			var skipNotify = false
			switch request.Op {
			case types.OpLogin:
				message = fmt.Sprintf("Version: %v, HostName: %v, Os: %v, Arch: %v",
					content["version"], content["hostname"], content["os"], content["arch"])
			case types.OpNewProxy:
				message = fmt.Sprintf("ProxyName: %v, ProxyType: %v, RemotePort: %v",
					content["proxy_name"], content["proxy_type"], content["remote_port"])
			case types.OpNewWorkConn:
				message = fmt.Sprintf("RunID: %v", content["run_id"])
			case types.OpNewUserConn:
				// ip cache
				remoteIP := strings.Split(fmt.Sprint(content["remote_addr"]), ":")[0]
				proxyName := fmt.Sprint(content["proxy_name"])
				ipCache, ok := userConnCache[proxyName]
				if !ok {
					ipCache = make(map[string]bool)
					userConnCache[proxyName] = ipCache
				}
				if _, ok := ipCache[remoteIP]; ok {
					skipNotify = true
					break
				}
				ipCache[remoteIP] = false
				var ipCName string
				if m.addressService != nil {
					ipCName = m.addressService.Query(remoteIP)
				}
				if ipCName != "" {
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s, CName: %s", proxyName, content["proxy_type"], remoteIP, ipCName)
				} else {
					message = fmt.Sprintf("ProxyName: %s, ProxyType: %v, RemoteIP: %s", proxyName, content["proxy_type"], remoteIP)
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
