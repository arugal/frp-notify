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
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/controller"
	"github/arugal/frp-notify/pkg/types"
)

type whitelistHandler struct {
	whitelist []string
}

func NewWhitelistHandler() controller.HandlerChain {
	handler := &whitelistHandler{}
	config.RegisterConfigListener(handler.watchNotifyConfig())
	return handler
}

func (c whitelistHandler) Op(op string) bool {
	if op == types.OpNewUserConn {
		return true
	}
	return false
}

func (c whitelistHandler) Do(req *types.Request) (bool, *types.Response) {
	whitelist := c.whitelist
	if len(whitelist) > 0 {
		conn := req.Body.(*types.UserConn)

		allow := false
		for _, whiteIP := range whitelist {
			if whiteIP == conn.RemoteIP {
				allow = true
				break
			}
		}

		if !allow {
			log.Debugf("reject [%s][%s][%s], not on the whitelist.", conn.ProxyType, conn.ProxyName, conn.RemoteIP)
			return false, &types.Response{
				Reject:       true,
				RejectReason: "remote ip is not on the whitelist",
				UnChange:     false,
				Content:      nil,
			}
		}
	}
	return true, nil
}

func (c *whitelistHandler) watchNotifyConfig() config.WatchNotifyConfigFunc {
	return func(cfg config.FRPNotifyConfig) {
		c.whitelist = cfg.Whitelist
	}
}
