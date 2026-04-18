// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package validator

import "regexp"

var (
	UsernameRe = regexp.MustCompile(`^[a-zA-Z0-9+-]{1,26}$`)
	SlugRe     = regexp.MustCompile(`^[a-zA-Z0-9-]{1,50}$`)
)

var (
	PasswordRe = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{}:;,.?]{6,18}$`)
)
