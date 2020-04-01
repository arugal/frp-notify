package models

const (
	APIVersion = "0.1.0"

	OpLogin       = "Login"
	OpNewProxy    = "NewProxy"
	OpNewAccessIp = "TraceAccessIp"
)

type Request struct {
	Version string      `json:"version"`
	Op      string      `json:"op"`
	Content interface{} `json:"content"`
}

type Response struct {
	Reject       bool        `json:"reject"`
	RejectReason string      `json:"reject_reason"`
	Unchange     bool        `json:"unchange"`
	Content      interface{} `json:"content"`
}

type LoginContent struct {
	Login
}

type UserInfo struct {
	User  string            `json:"user"`
	Metas map[string]string `json:"metas"`
}

type NewProxyContent struct {
	User UserInfo `json:"user"`
	NewProxy
}

type NewAccessIpContent struct {
	ProxyName    string `json:"proxy_name"`
	UserRemoteIp string `json:"user_remote_ip"`
}

// When frpc start, client send this message to login to server.
type Login struct {
	Version      string            `json:"version"`
	Hostname     string            `json:"hostname"`
	Os           string            `json:"os"`
	Arch         string            `json:"arch"`
	User         string            `json:"user"`
	PrivilegeKey string            `json:"privilege_key"`
	Timestamp    int64             `json:"timestamp"`
	RunId        string            `json:"run_id"`
	Metas        map[string]string `json:"metas"`

	// Some global configures.
	PoolCount int `json:"pool_count"`
}

// When frpc login success, send this message to frps for running a new proxy.
type NewProxy struct {
	ProxyName      string            `json:"proxy_name"`
	ProxyType      string            `json:"proxy_type"`
	UseEncryption  bool              `json:"use_encryption"`
	UseCompression bool              `json:"use_compression"`
	Group          string            `json:"group"`
	GroupKey       string            `json:"group_key"`
	Metas          map[string]string `json:"metas"`

	// tcp and udp only
	RemotePort int `json:"remote_port"`

	// http and https only
	CustomDomains     []string          `json:"custom_domains"`
	SubDomain         string            `json:"subdomain"`
	Locations         []string          `json:"locations"`
	HttpUser          string            `json:"http_user"`
	HttpPwd           string            `json:"http_pwd"`
	HostHeaderRewrite string            `json:"host_header_rewrite"`
	Headers           map[string]string `json:"headers"`

	// stcp
	Sk string `json:"sk"`

	// tcpmux
	Multiplexer string `json:"multiplexer"`
}
