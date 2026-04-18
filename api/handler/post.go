// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"code.cn/blog/api/common"
	"code.cn/blog/api/common/code"
	"code.cn/blog/api/middleware"
	"code.cn/blog/internal/dto/req"
	"code.cn/blog/internal/dto/resp"
	"code.cn/blog/internal/service"
	"code.cn/blog/pkg/validator"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	srv *service.PostService
}

func NewPostHandler(srv *service.PostService) *PostHandler {
	return &PostHandler{srv}
}

func (h *PostHandler) WebList(c *gin.Context) {
	var param req.PostWebList
	if err := c.ShouldBindQuery(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	ctx := c.Request.Context()
	list, count, err := h.srv.WebList(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OkData(c.Writer, resp.ListResp{
		Count:   count,
		Content: list,
	})
}

func (h *PostHandler) WebDetail(c *gin.Context) {
	var param req.PostWebDetailPathParams
	if err := c.ShouldBindUri(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	ctx := c.Request.Context()
	info, err := h.srv.WebDetail(ctx, param.Slug)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OkData(c.Writer, info)
}

func (h *PostHandler) List(c *gin.Context) {
	var param req.PostList
	if err := c.ShouldBindQuery(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	list, count, err := h.srv.List(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OkData(c.Writer, resp.ListResp{
		Count:   count,
		Content: list,
	})
}

func (h *PostHandler) Info(c *gin.Context) {
	var param req.InfoPathParams
	if err := c.ShouldBindUri(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	info, err := h.srv.Info(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OkData(c.Writer, info)
}

func (h *PostHandler) Create(c *gin.Context) {
	var param req.PostCreate
	if err := c.ShouldBindJSON(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	if !validator.SlugRe.MatchString(param.Slug) {
		common.FailMessage(c.Writer, "invalid slug")
		return
	}

	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	err := h.srv.Create(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OK(c.Writer)
}

func (h *PostHandler) Update(c *gin.Context) {
	var pathParam req.PostUpdatePathParams
	if err := c.ShouldBindUri(&pathParam); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	var param req.PostUpdate
	if err := c.ShouldBindJSON(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	if !validator.SlugRe.MatchString(param.Slug) {
		common.FailMessage(c.Writer, "invalid slug")
		return
	}

	param.ID = pathParam.ID
	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	err := h.srv.Update(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OK(c.Writer)
}

func (h *PostHandler) Delete(c *gin.Context) {
	var param req.DeletePathParams
	if err := c.ShouldBindUri(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	err := h.srv.Delete(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OK(c.Writer)
}

func (h *PostHandler) BatchDelete(c *gin.Context) {
	var param req.BatchDelete
	if err := c.ShouldBindJSON(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	param.UserID = middleware.GetUserID(c)

	ctx := c.Request.Context()
	err := h.srv.BatchDelete(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OK(c.Writer)
}
