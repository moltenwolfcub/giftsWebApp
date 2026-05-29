package config

const defaultPort = 8040

type Config struct {
	Port int
}

func New() *Config {
	return &Config{
		Port: defaultPort,
	}
}
