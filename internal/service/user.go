// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package service

import "code.cn/blog/internal/repository"

type UserService struct {
	repo          *repository.UserRepository
	userTokenRepo *repository.UserTokenRepository
}

func NewUserService(
	repo *repository.UserRepository,
	userTokenRepo *repository.UserTokenRepository,
) *UserService {
	return &UserService{repo, userTokenRepo}
}
