// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ctxUserID    = "uid"
	ctxSessionID = "sid"
)

func GetUserID(c *gin.Context) int {
	return c.GetInt(ctxUserID)
}

func GetSessionID(c *gin.Context) uuid.UUID {
	val, exists := c.Get(ctxSessionID)
	if !exists {
		return uuid.Nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}
