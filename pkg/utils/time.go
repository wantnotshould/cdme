// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package utils

import "time"

var (
	CNLoc     = time.FixedZone("CST", 8*3600)
	UnixEpoch = time.Unix(0, 0).UTC()
)

func Now() time.Time {
	return time.Now().In(CNLoc)
}

func ToCN(t time.Time) time.Time {
	return t.In(CNLoc)
}
