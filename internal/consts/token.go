// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package consts

import "time"

const (
	ATName     = "acctok"
	ATDuration = 1 * time.Hour
	ATMaxAge   = 3600

	RTName     = "reftok"
	RTDuration = 7 * 24 * time.Hour
	RTMaxAge   = 7 * 24 * 3600
)
