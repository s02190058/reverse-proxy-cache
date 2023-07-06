package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		GRPC    GRPC    `yaml:"grpc" env-prefix:"GRPC_"`
		Redis   Redis   `yaml:"redis" env-prefix:"REDIS_"`
		Youtube Youtube `yaml:"youtube" env-prefix:"YOUTUBE_"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"PORT"`
	}

	Redis struct {
		Host        string        `env-required:"true" yaml:"host" env:"HOST"`
		Port        string        `env-required:"true" yaml:"port" env:"PORT"`
		Password    string        `env-required:"true" env:"PASSWORD"`
		DB          int           `env-required:"true" yaml:"db" env:"DB"`
		DialTimeout time.Duration `env-required:"true" yaml:"dial_timeout" env:"DIAL_TIMEOUT"`
		ConnTimeout time.Duration `env-required:"true" yaml:"conn_timeout" env:"CONN_TIMEOUT"`
		TTL         time.Duration `env-required:"true" yaml:"ttl" env:"TTL"`
	}

	Youtube struct {
		Scheme  string        `env-required:"true" yaml:"scheme" env:"SCHEME"`
		Host    string        `env-required:"true" yaml:"host" env:"HOST"`
		Timeout time.Duration `env-required:"true" yaml:"timeout" env:"TIMEOUT"`
	}
)

func New(path string) (*Config, error) {
	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
