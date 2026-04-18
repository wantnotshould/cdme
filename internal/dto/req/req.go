// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package req

type Pagination struct {
	Page    int `form:"page" binding:"required,min=1"`
	PerPage int `form:"per_page" binding:"required,oneof=10 20 50"`
}

type InfoPathParams struct {
	InfoID int `uri:"id" binding:"required"`
	UserID int `json:"-"`
}

type DeletePathParams struct {
	DeleteID int `uri:"id" binding:"required"`
	UserID   int `json:"-"`
}

type BatchDelete struct {
	DeleteIds []int `json:"ids" binding:"required,min=1,dive,min=1"`
	UserID    int   `json:"-"`
}
