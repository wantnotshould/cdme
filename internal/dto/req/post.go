// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package req

type PostWebList struct {
	Pagination
}

type PostWebDetailPathParams struct {
	Slug string `uri:"slug" binding:"required,max=50"`
}

type PostList struct {
	Keyword    string `form:"keyword,omitempty" binding:"max=26"`
	Status     uint8  `form:"status,omitempty"`
	CategoryID uint8  `form:"category_id,omitempty"`
	Pagination
	UserID int `json:"-"`
}

type PostCreate struct {
	CategoryID uint8  `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=100"`
	Summary    string `json:"summary" binding:"required,max=200"`
	Slug       string `json:"slug" binding:"required,max=50"`
	Content    string `json:"content" binding:"required"`
	Status     uint8  `json:"status" binding:"required"`
	UserID     int    `json:"-"`
}

type PostUpdatePathParams struct {
	ID int `uri:"id" binding:"required"`
}

type PostUpdate struct {
	CategoryID uint8  `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=100"`
	Summary    string `json:"summary" binding:"required,max=200"`
	Slug       string `json:"slug" binding:"required,max=50"`
	Content    string `json:"content" binding:"required"`
	Status     uint8  `json:"status" binding:"required"`
	ID         int    `json:"-"`
	UserID     int    `json:"-"`
}
