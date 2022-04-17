package config

var Conf *Config

type Config struct {
	HttpReceiver bool // 启用http接收器
}

func NewConfig() *Config {
	conf := &Config{}
	Conf = conf
	return conf
}

func (conf *Config) Check() error {
	return nil
}
