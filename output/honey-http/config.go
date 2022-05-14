package http

const (
	// 默认push地址
	defPushAddress = "http://127.0.0.1:8080/push"
	// 默认压缩器名
	defCompress = "zstd"
	// 默认序列化器名
	defaultSerializer = "msgpack"
	// 默认请求超时
	defReqTimeout = 5
)

type Config struct {
	Disable      bool   // 关闭
	PushAddress  string // push地址, 示例: http://127.0.0.1:8080/push
	Compress     string // 压缩器名
	Serializer   string // 序列化器名
	AuthToken    string // 验证token, 如何设置, 请求header必须带上 token={AuthToken}, 如 token=myAuthToken
	ReqTimeout   int    // 请求超时, 单位秒
	ProxyAddress string // 代理地址. 支持 http, https, socks5, socks5h. 示例: socks5://127.0.0.1:1080
	ProxyUser    string // 代理用户名
	ProxyPasswd  string // 代理用户密码
}

func newConfig() *Config {
	return &Config{}
}

func (c *Config) Check() error {
	if c.PushAddress == "" {
		c.PushAddress = defPushAddress
	}
	if c.Compress == "" {
		c.Compress = defCompress
	}
	if c.Serializer == "" {
		c.Serializer = defaultSerializer
	}
	if c.ReqTimeout < 1 {
		c.ReqTimeout = defReqTimeout
	}
	return nil
}
