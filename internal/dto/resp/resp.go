// Copyright ©2026 cdme. All rights reserved.
// Author: cdme <https://cdme.cn>
// Email: hi@cdme.cn

package resp

type ListResp struct {
	Count   int64 `json:"count"`
	Content any `json:"content"`
}
