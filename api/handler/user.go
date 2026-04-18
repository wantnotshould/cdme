// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"code.cn/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{srv}
}

func (h *UserHandler) Login(c *gin.Context) {}

func (h *UserHandler) Profile(c *gin.Context) {}

func (h *UserHandler) RefreshToken(c *gin.Context) {}

func (h *UserHandler) Logout(c *gin.Context) {}
