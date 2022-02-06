package config

import (
	"emperror.dev/errors"
	"github.com/rome314/3-commas-ws-forwarder/internal/sender"
	"github.com/rome314/3-commas-ws-forwarder/pkg/connections"
)

var cfg *Config

type Config struct {
	Debug bool
	Redis connections.RedisConfig
	App   sender.Config
	Api   ApiConfig
}

type ApiConfig struct {
	Key    string
	Secret string
}

func (a ApiConfig) Valid() error {
	if a.Key == "" {
		return errors.New("key not provided")
	}
	if a.Secret == "" {
		return errors.New("secret not provided")

	}
	return nil
}

func GetConfig() *Config {
	return cfg
}

func (c *Config) Valid() error {
	if err := c.Redis.Valid(); err != nil {
		return errors.WithMessage(err, "redis")
	}
	if err := c.App.Valid(); err != nil {
		return errors.WithMessage(err, "app")
	}
	if err := c.Api.Valid(); err != nil {
		return errors.WithMessage(err, "api")
	}
	return nil

}
