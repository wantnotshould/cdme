// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"code.cn/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	srv *service.PostService
}

func NewPostHandler(srv *service.PostService) *PostHandler {
	return &PostHandler{srv}
}

func (h *PostHandler) WebList(c *gin.Context) {}

func (h *PostHandler) WebDetail(c *gin.Context) {}

func (h *PostHandler) List(c *gin.Context) {}

func (h *PostHandler) Info(c *gin.Context) {}

func (h *PostHandler) Create(c *gin.Context) {}

func (h *PostHandler) Update(c *gin.Context) {}

func (h *PostHandler) Delete(c *gin.Context) {}

func (h *PostHandler) BatchDelete(c *gin.Context) {}
