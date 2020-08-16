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
	"github/arugal/frp-notify/pkg/types"
	"testing"
)

func Test_whitelistHandler_Op(t *testing.T) {
	type args struct {
		op   string
		want bool
	}
	tests := []struct {
		name string
		args []args
	}{
		{
			"op",
			[]args{
				{
					op:   types.OpLogin,
					want: false,
				},
				{
					op:   types.OpNewProxy,
					want: false,
				},
				{
					op:   types.OpPing,
					want: false,
				},
				{
					op:   types.OpNewWorkConn,
					want: false,
				},
				{
					op:   types.OpNewUserConn,
					want: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := whitelistHandler{}
			for _, args := range tt.args {
				if got := n.Op(args.op); got != args.want {
					t.Errorf("Op() = %v, want %v", got, args.want)
				}
			}

		})
	}
}

func Test_whitelistHandler_Do(t *testing.T) {
	type fields struct {
		whitelist []string
	}
	type args struct {
		req *types.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"allow",
			fields{
				whitelist: []string{"127.0.0.1"},
			},
			args{
				req: &types.Request{
					Body: &types.UserConn{
						RemoteIP:  "127.0.0.1",
						ProxyType: "tcp",
						ProxyName: "allow",
					},
				},
			},
			true,
		},
		{
			"deny",
			fields{
				whitelist: []string{"127.0.0.1"},
			},
			args{
				req: &types.Request{
					Body: &types.UserConn{
						RemoteIP:  "127.0.0.2",
						ProxyType: "tcp",
						ProxyName: "deny",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := whitelistHandler{
				whitelist: tt.fields.whitelist,
			}
			got, _ := c.Do(tt.args.req)
			if got != tt.want {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}
