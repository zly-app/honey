package config

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (conf *Config) Check() error {
	return nil
}
