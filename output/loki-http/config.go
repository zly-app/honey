package http

const (
	// 默认push地址
	defPushAddress = "http://127.0.0.1:3100/loki/api/v1/push"
	// 默认启用压缩
	defEnableCompress = true
	// 默认请求超时
	defReqTimeout = 5
)

type Config struct {
	Disable        bool   // 关闭
	PushAddress    string // push地址, 示例: http://127.0.0.1:3100/loki/api/v1/push
	EnableCompress bool   // 是否启用压缩
	ReqTimeout     int    // 请求超时, 单位秒
}

func newConfig() *Config {
	return &Config{
		EnableCompress: defEnableCompress,
	}
}

func (c *Config) Check() error {
	if c.PushAddress == "" {
		c.PushAddress = defPushAddress
	}
	if c.ReqTimeout < 1 {
		c.ReqTimeout = defReqTimeout
	}
	return nil
}