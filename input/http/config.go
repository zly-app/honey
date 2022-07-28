package http

const (
	// 默认bind
	defaultBind = ":8080"
	// 默认推送路径
	defPushPath = "/push"
	// 默认post允许最大数据大小(32M)
	defaultPostMaxMemory = 32 << 20
)

type Config struct {
	Disable       bool   // 关闭
	Bind          string // 监听地址, 示例: :8080
	PushPath      string // 推送路径
	PostMaxMemory int64  // post允许客户端传输最大数据大小, 单位字节
	AuthToken     string // 验证token, 如果设置, 客户端请求header必须带上 token={AuthToken}, 如 token=myAuthToken
}

func newConfig() *Config {
	return &Config{}
}

func (c *Config) Check() error {
	if c.Bind == "" {
		c.Bind = defaultBind
	}
	if c.PushPath == "" {
		c.PushPath = defPushPath
	}
	if c.PostMaxMemory < 1 {
		c.PostMaxMemory = defaultPostMaxMemory
	}
	return nil
}
