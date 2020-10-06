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

func Test_notifyHandler_Op(t *testing.T) {
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
					want: true,
				},
				{
					op:   types.OpNewProxy,
					want: true,
				},
				{
					op:   types.OpPing,
					want: false,
				},
				{
					op:   types.OpNewWorkConn,
					want: true,
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
			n := notifyHandler{}
			for _, args := range tt.args {
				if got := n.Op(args.op); got != args.want {
					t.Errorf("Op() = %v, want %v", got, args.want)
				}
			}

		})
	}
}
