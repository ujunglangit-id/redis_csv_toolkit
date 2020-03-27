package util

import (
	"errors"
	"fmt"
	gcfg "gopkg.in/gcfg.v1"
)

type (
	// Config application configuration struct reflecting `.ini` file structure
	Config struct {
		RedisConfig RedisConfig
		AppConfig   AppConfig
	}

	AppConfig struct {
		KeyFormat    string
		FileLocation string
		FileName     string
	}

	// RedisConfig is a redis memeory cache configuration
	RedisConfig struct {
		Host      string
		MaxActive int
		MaxIdle   int
		Timeout   int
	}
)

func NewConfig() *Config {
	return &Config{}
}

// ReadConfig read `*.ini` configuration file and save it to variable of `*Config` type
func (cfg *Config) ReadConfig() error {
	var (
		e error
	)
	e = gcfg.ReadFileInto(cfg, "./config/toolkit_config.ini")
	if e != nil {
		return errors.New(fmt.Sprintf("failed to load file : %v\n", e))
	}

	return nil
}
