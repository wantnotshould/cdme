// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package handler

import (
	"errors"
	"net/http"

	"code.cn/blog/api/common"
	"code.cn/blog/api/common/code"
	"code.cn/blog/api/middleware"
	"code.cn/blog/conf"
	"code.cn/blog/internal/auth/token"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/consts"
	"code.cn/blog/internal/dto/req"
	"code.cn/blog/internal/service"
	"code.cn/blog/pkg/crypto/hash"
	"code.cn/blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{srv}
}

func (h *UserHandler) Login(c *gin.Context) {
	var param req.UserLogin
	if err := c.ShouldBind(&param); err != nil {
		common.Fail(c.Writer, code.ParamErr)
		return
	}

	param.IP = utils.IP(c.Request)
	param.UserAgent = utils.UserAgent(c.Request)

	ctx := c.Request.Context()
	res, err := h.srv.Login(ctx, param)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.SetAuthCookie(c.Writer, consts.ATName, res.AccessToken, consts.ATMaxAge)
	common.SetAuthCookie(c.Writer, consts.RTName, res.RefreshToken, consts.RTMaxAge)

	common.OK(c.Writer)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	ctx := c.Request.Context()
	info, err := h.srv.Profile(ctx, userID)
	if err != nil {
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.OkData(c.Writer, info)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(consts.RTName)
	if err != nil || refreshToken == "" {
		common.Custom(c.Writer, http.StatusUnauthorized, "session expired, please log in again")
		return
	}

	ip := utils.IP(c.Request)
	userAgent := utils.UserAgent(c.Request)

	claims, err := token.Parse(refreshToken, ip, userAgent)
	if err != nil {
		// If token parsing fails (invalid, expired, or tampered), clear cookies.
		common.CleanAuthCookie(c.Writer)
		message := "Invalid refresh token"
		if errors.Is(err, jwt.ErrTokenExpired) {
			message = "refresh token expired, please log in again"
		}

		common.Custom(c.Writer, http.StatusUnauthorized, message)
		return
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(utils.Now()) {
		common.CleanAuthCookie(c.Writer)
		common.Custom(c.Writer, http.StatusUnauthorized, "token expired")
		return
	}

	refreshParam := req.UserRefreshToken{
		RefreshToken: hash.HMACBlake2b256Hex([]byte(refreshToken), []byte(conf.Get().Hash.Key)),
		IP:           ip,
		UserAgent:    userAgent,
	}

	ctx := c.Request.Context()
	res, err := h.srv.RefreshToken(ctx, refreshParam, claims)
	if err != nil {
		common.CleanAuthCookie(c.Writer)
		common.FailMessage(c.Writer, err.Error())
		return
	}

	common.SetAuthCookie(c.Writer, consts.ATName, res.AccessToken, consts.ATMaxAge)
	common.SetAuthCookie(c.Writer, consts.RTName, res.RefreshToken, consts.RTMaxAge)

	common.OK(c.Writer)
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	if sessionID == uuid.Nil || userID == 0 {
		common.FailMessage(c.Writer, "session expired, please log in again")
		return
	}

	ctx := c.Request.Context()
	if err := h.srv.Logout(ctx, userID, sessionID); err != nil {
		common.FailMessage(c.Writer, "logout failed")
		return
	}

	redis.DB().DelAccessToken(ctx, userID)
	redis.DB().DelRefreshToken(ctx, userID)

	common.CleanAuthCookie(c.Writer)

	common.OK(c.Writer)
}
