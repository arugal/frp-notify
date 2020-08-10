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

type blacklistHandler struct {
	blacklist []string
}

func NewBlacklistHandler() controller.HandlerChain {
	handler := &blacklistHandler{}
	config.RegisterConfigListener(handler.watchNotifyConfig())
	return handler
}

func (b blacklistHandler) Op(op string) bool {
	if op == types.OpNewUserConn {
		return true
	}
	return false
}

func (b blacklistHandler) Do(req *types.Request) (bool, *types.Response) {
	blacklist := b.blacklist
	if len(blacklist) > 0 {
		conn := req.Body.(*types.UserConn)

		deny := false
		for _, blackIP := range blacklist {
			if blackIP == conn.RemoteIP {
				deny = true
				break
			}
		}

		if deny {
			log.Debugf("reject [%s][%s][%s], on the blacklist.", conn.ProxyType, conn.ProxyName, conn.RemoteIP)
			return false, &types.Response{
				Reject:       true,
				RejectReason: "remote ip is on the blacklist",
				UnChange:     false,
				Content:      nil,
			}
		}
	}
	return true, nil
}

func (b *blacklistHandler) watchNotifyConfig() config.WatchNotifyConfigFunc {
	return func(cfg config.FRPNotifyConfig) {
		b.blacklist = cfg.Blacklist
	}
}
