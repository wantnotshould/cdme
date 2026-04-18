// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package common

import (
	"net/http"

	"code.cn/blog/cmd/flags"
	"code.cn/blog/internal/consts"
)

func SetAuthCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   !flags.Debug,
		SameSite: http.SameSiteLaxMode,
	})
}

func CleanAuthCookie(w http.ResponseWriter) {
	SetAuthCookie(w, consts.ATName, "", -1)
	SetAuthCookie(w, consts.RTName, "", -1)
}
