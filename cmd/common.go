// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package cmd

import (
	"log"

	"code.cn/blog/conf"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/database"
	"code.cn/blog/pkg/crypto/aes"
)

func setup() {
	conf.Init()
	aes.Init([]byte(conf.Get().AESGCM.Key))
	redis.Init()

	database.Init()
	db := database.Get()
	if db == nil {
		log.Fatal("database not initialized")
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("database migrate failed:", err)
	}
}

func release() {
	redis.DB().Close()

	if database.Instance() != nil {
		_ = database.Instance().Close()
	}
}
