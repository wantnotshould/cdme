// Copyright ©2026 cdme. All rights reserved.
// Author: cdme <https://cdme.cn>
// Email: hi@cdme.cn

package resp

import "time"

type PostWebListItem struct {
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

type PostWebDetail struct {
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostListItem struct {
	ID         int       `json:"id"`
	CategoryID uint8     `json:"category_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Summary    string    `json:"summary"`
	Status     uint8     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostInfo struct {
	ID         int       `json:"id"`
	CategoryID uint8     `json:"category_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Summary    string    `json:"summary"`
	Content    string    `json:"content"`
	Status     uint8     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
