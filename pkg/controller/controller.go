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

package controller

import (
	"encoding/json"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/logger"
	"github/arugal/frp-notify/pkg/types"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

var normalResponse = types.Response{
	Reject:       false,
	RejectReason: "",
	UnChange:     true,
	Content:      nil,
}

var rejectedResponse = types.Response{
	Reject:       true,
	RejectReason: "",
	UnChange:     true,
	Content:      nil,
}

func init() {
	log = logger.Log
}

type ManagerController struct {
	chains []HandlerChain
	cfg    config.FRPNotifyConfig
}

func NewManagerController(opts ...ManagerControllerOption) *ManagerController {
	ms := &ManagerController{}

	config.RegisterConfigListener(func(cfg config.FRPNotifyConfig) {
		ms.cfg = cfg
	})

	for _, o := range opts {
		o(ms)
	}
	return ms
}

type ManagerControllerOption func(m *ManagerController)

func WithHandlerChains(handler ...HandlerChain) ManagerControllerOption {
	return func(m *ManagerController) {
		m.chains = handler
	}
}

// Start start manager service
func (m *ManagerController) Start(stop chan struct{}) error {
	for _, handler := range m.chains {
		if richHandler, ok := handler.(RichHandlerChain); ok {
			err := richHandler.Open()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Close
func (m *ManagerController) Close() {
	for _, handler := range m.chains {
		if richHandler, ok := handler.(RichHandlerChain); ok {
			_ = richHandler.Close()
		}
	}
}

func (m *ManagerController) Register(mux *http.ServeMux) {

	render := func(writer http.ResponseWriter, statusCode int, res interface{}) {
		writer.WriteHeader(statusCode)
		b, err := json.Marshal(res)
		if err != nil {
			log.Debugf("marshal err: %v", err)
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Debugf("response writer err: %v", err)
		}
	}

	mux.HandleFunc("/handler", func(writer http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("read request body error: %v \n", err)
			render(writer, http.StatusOK, rejectedResponse)
			return
		}

		request := &types.Request{}
		err = json.Unmarshal(body, request)
		if err != nil {
			log.Errorf("unmarshal request body error: %v \n", err)
			render(writer, http.StatusOK, rejectedResponse)
			return
		}

		// verify api version
		if request.Version != types.APIVersion {
			log.Warnf("unsupported api version %s \n", request.Version)
			render(writer, http.StatusOK, rejectedResponse)
			return
		}

		// ignore ping operator
		if request.Op == types.OpPing {
			log.Debug("ignore ping operation")
			render(writer, http.StatusOK, normalResponse)
			return
		}

		// parser content
		err = parserContent(request)
		if err != nil {
			log.Errorf("unmarshal request content [%s] error: %v \n", request.Content, err)
			render(writer, http.StatusOK, rejectedResponse)
			return
		}

		for _, handler := range m.chains {
			if handler.Op(request.Op) {
				ok, resp := handler.Do(request)
				if !ok {
					render(writer, http.StatusOK, &resp)
					return
				}
			}
		}
		render(writer, http.StatusOK, normalResponse)
	})
}

func parserContent(req *types.Request) error {
	var body interface{}
	switch req.Op {
	case types.OpLogin:
		body = &types.Login{}
	case types.OpNewUserConn:
		body = &types.UserConn{}
	case types.OpNewProxy:
		body = &types.Proxy{}
	case types.OpNewWorkConn:
		body = &types.WorkConn{}
	default:
		return nil
	}

	content, err := json.Marshal(req.Content)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, body)
	if err != nil {
		return err
	}

	if req.Op == types.OpNewUserConn {
		conn := body.(*types.UserConn)
		conn.RemoteIP, _, _ = net.SplitHostPort(conn.RemoteAddr)
	}

	req.Body = body
	return nil
}
