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

package types

const (
	APIVersion = "0.1.0"

	OpLogin       = "Login"
	OpNewProxy    = "NewProxy"
	OpCloseProxy  = "CloseProxy"
	OpPing        = "Ping"
	OpNewWorkConn = "NewWorkConn"
	OpNewUserConn = "NewUserConn"
)

type Request struct {
	Version string                 `json:"version"`
	Op      string                 `json:"op"`
	Content map[string]interface{} `json:"content,omitempty"`
	Body    interface{}            `json:"-"`
}

type Response struct {
	Content      interface{} `json:"content"`
	RejectReason string      `json:"reject_reason"`
	Reject       bool        `json:"reject"`
	UnChange     bool        `json:"unchange"`
}

type User struct {
	User  string            `json:"user"`
	Metas map[string]string `json:"metas,omitempty"`
	RunID string            `json:"run_id"`
}

type Login struct {
	Version      string            `json:"version"`
	Hostname     string            `json:"hostname"`
	OS           string            `json:"os"`
	Arch         string            `json:"arch"`
	User         string            `json:"user"`
	Timestamp    int64             `json:"timestamp"`
	PrivilegeKey string            `json:"privilege_key"`
	RunID        string            `json:"run_id"`
	PoolCount    int               `json:"pool_count"`
	Metas        map[string]string `json:"metas,omitempty"`
}

type Proxy struct {
	User           User   `json:"user"`
	ProxyName      string `json:"proxy_name"`
	ProxyType      string `json:"proxy_type"`
	UseEncryption  bool   `json:"use_encryption"`
	UseCompression bool   `json:"use_compression"`
	Group          string `json:"group"`
	GroupKey       string `json:"group_key"`

	RemotePort int `json:"remote_port"` // tcp and udp only

	CustomDomains     []string          `json:"custom_domains"` // http and https only
	Subdomain         string            `json:"subdomain"`
	Locations         []string          `json:"locations"`
	HTTPUser          string            `json:"http_user"`
	HTTPPwd           string            `json:"http_pwd"`
	HostHeaderRewrite string            `json:"host_header_rewrite"`
	Headers           map[string]string `json:"headers,omitempty"`

	SK string `json:"sk"` // stcp only

	Multiplexer string `json:"multiplexer"` // tcpmux only

	Metas map[string]string `json:"metas,omitempty"`
}

type CloseProxy struct {
	User      User   `json:"user"`
	ProxyName string `json:"proxy_name"`
}

type WorkConn struct {
	User         User   `json:"user"`
	RunID        string `json:"run_id"`
	Timestamp    string `json:"timestamp"`
	PrivilegeKey string `json:"privilege_key"`
}

type UserConn struct {
	User       User   `json:"user"`
	ProxyName  string `json:"proxy_name"`
	ProxyType  string `json:"proxy_type"`
	RemoteAddr string `json:"remote_addr"`
	RemoteIP   string `json:"-"`
}
