// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package cmd

import (
	"code.cn/blog/conf"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/pkg/crypto/aes"
)

func setup() {
	conf.Init()
	aes.Init([]byte(conf.Get().AESGCM.Key))
	redis.Init()
}

func release() {
	redis.DB().Close()
}
