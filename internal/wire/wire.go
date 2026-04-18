//go:build wireinject
// +build wireinject

// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package wire

import (
	"code.cn/blog/api/handler"
	"code.cn/blog/internal/repository"
	"code.cn/blog/internal/service"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type App struct {
	UserHandler *handler.UserHandler
	PostHandler *handler.PostHandler
}

var providerSet = wire.NewSet(
	handler.NewUserHandler,
	service.NewUserService,
	repository.NewUserRepository,
	repository.NewUserTokenRepository,
	handler.NewPostHandler,
	service.NewPostService,
	repository.NewPostRepository,
	wire.Struct(new(App), "*"),
)

func Init(db *gorm.DB) *App {
	wire.Build(providerSet)
	return &App{}
}
