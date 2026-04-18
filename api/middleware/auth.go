// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package middleware

import (
	"code.cn/blog/api/common"
	"code.cn/blog/internal/auth/token"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/consts"
	"code.cn/blog/pkg/crypto/hash"
	"code.cn/blog/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Anonymous() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken, err := c.Cookie(consts.ATName)
		if err != nil || accessToken == "" {
			common.Unauthorized(c.Writer)
			c.Abort()
			return
		}

		claims, err := token.Parse(
			accessToken,
			utils.IP(c.Request),
			utils.UserAgent(c.Request),
		)

		if err != nil {
			common.CleanAuthCookie(c.Writer)
			common.Unauthorized(c.Writer)
			c.Abort()
			return
		}

		accessTokenHash := hash.HashBlake2b256Hex([]byte(accessToken))

		ctx := c.Request.Context()
		storedHash, err := redis.DB().GetAccessToken(
			ctx,
			claims.DecryptedPayload.UserID,
		)

		if err != nil {
			common.Unauthorized(c.Writer)
			c.Abort()
			return
		}

		if storedHash == "" || accessTokenHash != storedHash {
			common.Unauthorized(c.Writer)
			c.Abort()
			return
		}

		c.Set(ctxSessionID, claims.SessionID)
		c.Set(ctxUserID, claims.DecryptedPayload.UserID)

		c.Next()
	}
}
