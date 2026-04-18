// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package router

import (
	"code.cn/blog/api/common"
	"code.cn/blog/api/middleware"
	"code.cn/blog/internal/database"
	"code.cn/blog/internal/wire"
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	e.Use(middleware.Anonymous())

	db := database.Get()
	app := wire.Init(db)
	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", app.UserHandler.Login)
		auth.GET("/refresh-token", app.UserHandler.RefreshToken)
	}

	admin := v1.Group("/admin")
	admin.Use(middleware.Auth())
	{
		user := admin.Group("/users")
		{
			user.GET("/profile", app.UserHandler.Profile)
			user.POST("/logout", app.UserHandler.Logout)
		}

		post := admin.Group("/posts")
		{
			post.GET("", app.PostHandler.List)
			post.GET("/:id", app.PostHandler.Info)
			post.POST("", app.PostHandler.Create)
			post.PUT("/:id", app.PostHandler.Update)
			post.DELETE("/:id", app.PostHandler.Delete)
			post.POST("/batch-delete", app.PostHandler.BatchDelete)
		}
	}

	// Web routes
	site := v1.Group("/site")
	{
		post := site.Group("/posts")
		{
			post.GET("", app.PostHandler.WebList)
			post.GET("/:slug", app.PostHandler.WebDetail)
		}
	}

	e.NoRoute(func(c *gin.Context) {
		common.NotFound(c.Writer)
	})
}
