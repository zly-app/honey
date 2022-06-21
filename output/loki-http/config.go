package http

const (
	// 默认push地址
	defPushAddress = "http://127.0.0.1:3100/loki/api/v1/push"
	// 默认启用压缩
	defEnableCompress = true
	// 默认请求超时
	defReqTimeout = 5
	// 默认请求重试次数
	defRetryCount = 2
	// 默认请求重试间隔时间
	defRetryIntervalMs = 2000
)

type Config struct {
	Disable         bool   // 关闭
	PushAddress     string // push地址, 示例: http://127.0.0.1:3100/loki/api/v1/push
	EnableCompress  bool   // 是否启用压缩
	ReqTimeout      int    // 请求超时, 单位秒
	RetryCount      int    // 请求失败重试次数, 0表示禁用
	RetryIntervalMs int    // 请求失败重试间隔毫秒数
	ProxyAddress    string // 代理地址. 支持 http, https, socks5, socks5h. 示例: socks5://127.0.0.1:1080 socks5://user:pwd@127.0.0.1:1080
}

func newConfig() *Config {
	return &Config{
		EnableCompress:  defEnableCompress,
		RetryCount:      defRetryCount,
		RetryIntervalMs: defRetryIntervalMs,
	}
}

func (c *Config) Check() error {
	if c.PushAddress == "" {
		c.PushAddress = defPushAddress
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
