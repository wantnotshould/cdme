// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package middleware

import "github.com/gin-gonic/gin"

func Anonymous() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
