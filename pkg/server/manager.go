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
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/types"
	"io/ioutil"
	"net/http"
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

type managerServer struct {
	serverAddr     string
	windowInterval time.Duration
	requestChan    chan *types.Request
}

func NewManagerServer(opts ...ManagerServerOption) *managerServer {
	ms := &managerServer{
		serverAddr:     defaultServerAddr,
		windowInterval: defaultWindowInterval,
		requestChan:    make(chan *types.Request, maxRequestSize),
	}

	for _, o := range opts {
		o(ms)
	}
	return ms
}

type ManagerServerOption func(m *managerServer)

func WithWindowInterval(windowInterval time.Duration) ManagerServerOption {
	return func(m *managerServer) {
		m.windowInterval = windowInterval
	}
}

func WithServerAddr(serverAddr string) ManagerServerOption {
	return func(m *managerServer) {
		m.serverAddr = serverAddr
	}
}

func (m *managerServer) HttpServer() {
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
		m.requestChan <- request
	})

	err := r.Run(m.serverAddr)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *managerServer) doNotify() {
	for request := range m.requestChan {
		log.Debug("receive new request op: %s, content: %s \n", request.Op, request.Content)

		content, ok := request.Content.(map[string]interface{})
		if !ok {
			log.Warnf("unsupported content: %s \n", request.Content)
			return
		}
		var title, message string
		switch request.Op {
		case types.OpLogin:
		case types.OpNewProxy:
		case types.OpPing:
		case types.OpNewWorkConn:
		case types.OpNewUserConn:
		}
		log.Infof("notify: %s  %s \n", title, message)
	}
}

func (m *managerServer) Close() {
	close(m.requestChan)
}
