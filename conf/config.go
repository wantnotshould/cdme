// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package conf

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"code.cn/blog/cmd/flags"
	"code.cn/blog/pkg/utils"
)

type Scheme struct {
	Port string `json:"port"`
}

type Config struct {
	Scheme Scheme `json:"scheme"`
}

func (c *Config) validate() error {
	if c.Scheme.Port == "" || !strings.HasPrefix(c.Scheme.Port, ":") {
		return utils.Err("scheme.Port is empty or format error")
	}
	return nil
}

var (
	cfgPtr   atomic.Pointer[Config]
	fullPath string
	once     sync.Once
)

func Get() *Config {
	return cfgPtr.Load()
}

func defaultConfig() *Config {
	return &Config{
		Scheme: Scheme{
			Port: ":2603",
		},
	}
}

func load() error {
	if fullPath == "" {
		return utils.Err("config path not initialized, call Init() first")
	}

	var newConfig Config
	if err := utils.ReadJSON(fullPath, &newConfig); err != nil {
		return utils.Wrap("failed to load config file", err)
	}

	if err := newConfig.validate(); err != nil {
		return utils.Wrap("config validation failed", err)
	}

	cfgPtr.Store(&newConfig)

	return nil
}

func Init() {
	once.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get working directory: %v\n", err)
		}

		fullPath = filepath.Join(wd, flags.Data, "config.json")

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			def := defaultConfig()
			if err := utils.WriteJSON(fullPath, &def); err != nil {
				log.Fatalf("failed to initialize config file: %v", err)
			}
			cfgPtr.Store(def)
		} else {
			if err := load(); err != nil {
				log.Fatalf("failed to load config file: %v", err)
			}
		}
	})
}
