// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package middleware

import (
	"bytes"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cn/blog/api/common"
	"code.cn/blog/conf"
	"code.cn/blog/internal/auth/token"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/consts"
	"code.cn/blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/xiayoudi/ud/hash"
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func absTimeDiff(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}

func Anonymous() gin.HandlerFunc {
	return func(c *gin.Context) {

		tsStr := c.GetHeader("X-Timestamp")
		if tsStr == "" {
			common.Custom(c.Writer, http.StatusUnauthorized, "bad request")
			c.Abort()
			return
		}

		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil || absTimeDiff(time.Now().Unix(), ts) > 10 {
			common.Custom(c.Writer, http.StatusUnauthorized, "expired")
			c.Abort()
			return
		}

		clientSig := c.GetHeader("X-Signature")
		if clientSig == "" {
			common.Custom(c.Writer, http.StatusUnauthorized, "bad request")
			c.Abort()
			return
		}

		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)

		// base
		buf.WriteString(c.Request.URL.Path)
		buf.WriteString(tsStr)
		buf.WriteString(c.Request.Method)

		// query (stable + full values)
		if len(c.Request.URL.Query()) > 0 {
			keys := make([]string, 0, len(c.Request.URL.Query()))
			for k := range c.Request.URL.Query() {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				values := c.Request.URL.Query()[k]
				for _, v := range values {
					buf.WriteString(k)
					buf.WriteString(v)
				}
			}
		}

		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		// body handling
		if !isMultipart && c.Request.Body != nil {

			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				common.Custom(c.Writer, http.StatusRequestEntityTooLarge, "payload too large")
				c.Abort()
				return
			}

			// restore body
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			buf.Write(bodyBytes)
		}

		// avoid string() copy if possible
		expectedSig := hash.HMACBlake2b256Hex(buf.Bytes(), []byte(conf.Get().Scheme.PublicKey))

		if expectedSig != clientSig {
			common.Custom(c.Writer, http.StatusUnauthorized, "sign error")
			c.Abort()
			return
		}

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

		accessTokenHash := hash.HMACBlake2b256Hex([]byte(accessToken), []byte(conf.Get().Hash.Key))

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
