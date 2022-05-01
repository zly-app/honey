package http

const (
	// 默认bind
	defaultBind = ":8080"
	// 默认上报路径
	defReportPath = "/report"
	// 默认post允许最大数据大小(32M)
	defaultPostMaxMemory = 32 << 20
	// 默认压缩器名
	defCompress = "zstd"
	// 默认序列化器名
	defaultSerializer = "msgpack"
)

type Config struct {
	Disable       bool   // 关闭
	Bind          string // 监听地址, 示例: :8080
	ReportPath    string // 上报路径
	PostMaxMemory int64  // post允许客户端传输最大数据大小, 单位字节
	Compress      string // 压缩器名
	Serializer    string // 序列化器名
	AuthToken     string // 验证token, 如何设置, 请求header必须带上 token={AuthToken}, 如 token=myAuthToken
}

func newConfig() *Config {
	return &Config{}
}

func (c *Config) Check() error {
	if c.Bind == "" {
		c.Bind = defaultBind
	}
	if c.ReportPath == "" {
		c.ReportPath = defReportPath
	}
	if c.PostMaxMemory < 1 {
		c.PostMaxMemory = defaultPostMaxMemory
	}
	if c.Compress == "" {
		c.Compress = defCompress
	}
	if c.Serializer == "" {
		c.Serializer = defaultSerializer
	}
	return nil
}
