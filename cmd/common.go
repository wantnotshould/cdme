// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package cmd

import (
	"context"
	"log"
	"sync"
	"time"

	"code.cn/blog/conf"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/database"
	"code.cn/blog/internal/logger"
	"code.cn/blog/internal/repository"
	"github.com/xiayoudi/ud/aes"
)

var (
	cleanupCancel context.CancelFunc
	cleanupWG     sync.WaitGroup
)

func runCleanupLoop(
	ctx context.Context,
	job func(context.Context) (int64, error),
	interval time.Duration,
) {
	defer cleanupWG.Done()

	_, _ = job(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, _ = job(ctx)

		case <-ctx.Done():
			return
		}
	}
}

func startCleanupTasks() {
	var ctx context.Context
	ctx, cleanupCancel = context.WithCancel(context.Background())

	repo := repository.NewUserTokenRepository(database.Get())

	cleanupWG.Add(1)
	go runCleanupLoop(ctx, repo.CleanupExpired, 1*time.Hour)

	cleanupWG.Add(1)
	go runCleanupLoop(ctx, repo.CleanupRevoked, 6*time.Hour)
}

func setup() {
	conf.Init()
	logger.Init()
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

	startCleanupTasks()
}

func release() {
	if cleanupCancel != nil {
		cleanupCancel()
	}

	cleanupWG.Wait()

	redis.DB().Close()

	if database.Instance() != nil {
		_ = database.Instance().Close()
	}

	logger.Close()
}
