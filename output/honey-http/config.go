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
	// 默认请求重试次数
	defRetryCount = 2
	// 默认请求重试间隔时间
	defRetryIntervalMs = 2000
)

type Config struct {
	Disable         bool   // 关闭
	PushAddress     string // push地址, 示例: http://127.0.0.1:8080/push
	Compress        string // 压缩器名
	Serializer      string // 序列化器名
	AuthToken       string // 验证token, 如何设置, 请求header必须带上 token={AuthToken}, 如 token=myAuthToken
	ReqTimeout      int    // 请求超时, 单位秒
	RetryCount      int    // 请求失败重试次数, 0表示禁用
	RetryIntervalMs int    // 请求失败重试间隔毫秒数
	ProxyAddress    string // 代理地址. 支持 http, https, socks5, socks5h. 示例: socks5://127.0.0.1:1080 socks5://user:pwd@127.0.0.1:1080
}

func newConfig() *Config {
	return &Config{
		RetryCount:      defRetryCount,
		RetryIntervalMs: defRetryIntervalMs,
	}
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
	if c.RetryCount < 0 {
		c.RetryCount = 0
	}
	if c.RetryIntervalMs < 0 {
		c.RetryIntervalMs = 0
	}
	return nil
}
