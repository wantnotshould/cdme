// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package conf

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"code.cn/blog/cmd/flags"
	"github.com/xiayoudi/ud"
)

type Scheme struct {
	Port      string `json:"port"`
	PublicKey string `json:"public_key"`
}

type Hash struct {
	Key string `json:"key"`
}

type AESGCM struct {
	Key string `json:"key"`
	AAD string `json:"aad"`
}

type JWT struct {
	Key string `json:"key"`
}

type Redis struct {
	Addr         string        `json:"addr"`
	Password     string        `json:"password"`
	DB           int           `json:"db"`
	Prefix       string        `json:"prefix"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	MaxRetries   int           `json:"max_retries"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout"`
}

type Database struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	DBName       string        `json:"db_name"`
	Timezone     string        `json:"timezone"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

func (d *Database) DSN() string {
	tz := d.Timezone
	if tz == "" {
		tz = "Asia/Shanghai"
	}

	params := url.Values{}
	params.Set("charset", "utf8mb4")
	params.Set("parseTime", "True")
	params.Set("loc", tz)

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		params.Encode(),
	)
}

type Logger struct {
	LogsDir    string `json:"logs_dir"`
	MaxSize    int    `json:"max_size"` // MB
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
}

type Config struct {
	Scheme   Scheme   `json:"scheme"`
	Hash     Hash     `json:"hash"`
	AESGCM   AESGCM   `json:"aesgcm"`
	JWT      JWT      `json:"jwt"`
	Redis    Redis    `json:"redis"`
	Database Database `json:"database"`
	Logger   Logger   `json:"logger"`
}

func (c *Config) validate() error {
	if c.Scheme.Port == "" || !strings.HasPrefix(c.Scheme.Port, ":") {
		return ud.Err("scheme.Port is empty or format error")
	}

	if c.Scheme.PublicKey == "" {
		return ud.Err("scheme.public_key annot be empty")
	}

	allowedLens := []int{16, 24, 32}
	if !slices.Contains(allowedLens, len(c.Hash.Key)) {
		return ud.Err("invalid hash.key length (16/24/32)")
	}

	if !slices.Contains(allowedLens, len(c.AESGCM.Key)) {
		return ud.Err("invalid aesgcm.key length (16/24/32)")
	}

	if strings.TrimSpace(c.AESGCM.AAD) == "" {
		return ud.Err("aesgcm.aad cannot be empty (recommended for security)")
	}

	if len(c.JWT.Key) < 32 {
		return ud.Err("jwt.key too weak, use at least 32 bytes")
	}

	if c.Redis.Addr == "" {
		return ud.Err("redis address can't be empty")
	}

	if c.Redis.Prefix == "" {
		return ud.Err("redis prefix can't be empty")
	}

	if c.Redis.PoolSize <= 0 {
		return ud.Err("redis pool_size must be > 0")
	}

	if c.Redis.MinIdleConns <= 0 {
		return ud.Err("redis min_idle_conns must be > 0")
	}

	if c.Redis.MaxRetries < 0 {
		return ud.Err("redis max_retries must be >= 0")
	}

	if c.Redis.DialTimeout <= 0 {
		return ud.Err("redis dial_timeout must be > 0")
	}

	if c.Redis.ReadTimeout <= 0 {
		return ud.Err("redis read_timeout must be > 0")
	}

	if c.Redis.WriteTimeout <= 0 {
		return ud.Err("redis write_timeout must be > 0")
	}

	if c.Redis.PoolTimeout <= 0 {
		return ud.Err("redis pool_timeout must be > 0")
	}

	if c.Database.Host == "" {
		return ud.Err("database host can't be empty")
	}

	if c.Database.Port == "" {
		return ud.Err("database port can't be empty")
	}

	if c.Database.User == "" {
		return ud.Err("database user can't be empty")
	}

	if c.Database.Password == "" {
		return ud.Err("database password can't be empty")
	}

	if c.Database.DBName == "" {
		return ud.Err("database name can't be empty")
	}

	if c.Logger.LogsDir == "" {
		return ud.Err("logger logs_dir can't be empty")
	}

	if c.Logger.MaxSize <= 0 {
		return ud.Err("logger max_size must be greater than 0")
	}

	if c.Logger.MaxBackups < 0 {
		return ud.Err("logger max_backups can't be negative")
	}

	if c.Logger.MaxAge < 0 {
		return ud.Err("logger max_age can't be negative")
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
	logsDir := filepath.Join(flags.Data, "logs")

	return &Config{
		Scheme: Scheme{
			Port:      ":2603",
			PublicKey: "大佬别攻击，高抬贵手。感谢！ Guru, no hacking. thanks!",
		},
		Hash: Hash{
			// openssl rand -base64 32 | cut -c1-32
			Key: "ykSlmOR2yL9Et/lO4QeTgzDuvU0/GHVk",
		},
		AESGCM: AESGCM{
			// openssl rand -base64 32 | cut -c1-32
			Key: "57cVg1gFKk/zBavQeIGad7hbqe7MfUMf",
			// openssl rand -base64 32
			AAD: "u1Y+M42Y9R32oGcSAeHs7NZniyO7xLAG5tmMwW1h9ms=",
		},
		JWT: JWT{
			Key: "GO0nIDh1aPYK3Kzlv4Ljxwvta3F0aEKr8JOqHHsoVxQ=",
		},
		Redis: Redis{
			Addr:         "127.0.0.1:6379",
			Password:     "",
			DB:           0,
			Prefix:       "cdme_bog",
			PoolSize:     100,
			MinIdleConns: 10,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		},
		Database: Database{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Password:     "root",
			DBName:       "blog",
			Timezone:     "Asia/Shanghai",
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			MaxLifetime:  time.Hour,
		},
		Logger: Logger{
			LogsDir:    logsDir,
			MaxSize:    50,
			MaxBackups: 10,
			MaxAge:     24,
		},
	}
}

func load() error {
	if fullPath == "" {
		return ud.Err("config path not initialized, call Init() first")
	}

	var newConfig Config
	if err := ud.ReadJSON(fullPath, &newConfig); err != nil {
		return ud.Wrap(err, "failed to load config file")
	}

	if err := newConfig.validate(); err != nil {
		return ud.Wrap(err, "config validation failed")
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
			if err := ud.WriteJSON(fullPath, &def); err != nil {
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
