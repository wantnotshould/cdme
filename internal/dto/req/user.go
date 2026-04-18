// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package req

type UserLogin struct {
	Username  string `json:"username" binding:"required,min=1,max=26"`
	Password  string `json:"password" binding:"required,min=6,max=18"`
	IP        string `json:"-"`
	UserAgent string `json:"-"`
}

type UserRefreshToken struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
	IP           string `json:"-"`
	UserAgent    string `json:"-"`
}

type UserLogout struct {
	AccessToken string `form:"access_token" binding:"required"`
	IP          string `json:"-"`
	UserAgent   string `json:"-"`
}
