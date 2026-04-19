// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package model

const (
	PostStatusDraft uint8 = iota + 1
	PostStatusPublished
	PostStatusUnderReview
	PostStatusPrivate
)

const (
	PostCategorySharing uint8 = iota + 1
	PostCategoryTech
	PostCategoryEssay
	PostCategoryTutorial
	PostCategoryReading
)

type PostOption struct {
	ID    uint8  `json:"id"`
	Label string `json:"label"`
}

var PostStatusOptions = []PostOption{
	{ID: PostStatusDraft, Label: "草稿"},
	{ID: PostStatusPublished, Label: "已发布"},
	{ID: PostStatusUnderReview, Label: "待审核"},
	{ID: PostStatusPrivate, Label: "私有"},
}

var PostCategoryOptions = []PostOption{
	{ID: PostCategorySharing, Label: "分享"},
	{ID: PostCategoryTech, Label: "技术"},
	{ID: PostCategoryEssay, Label: "随笔"},
	{ID: PostCategoryTutorial, Label: "教程"},
	{ID: PostCategoryReading, Label: "读书笔记"},
}

func toMap(opts []PostOption) map[uint8]PostOption {
	m := make(map[uint8]PostOption, len(opts))
	for _, o := range opts {
		m[o.ID] = o
	}
	return m
}

var (
	PostStatusMap   = toMap(PostStatusOptions)
	PostCategoryMap = toMap(PostCategoryOptions)
)

func IsValidPostStatus(id uint8) bool {
	_, ok := PostStatusMap[id]
	return ok
}

func IsValidPostCategory(id uint8) bool {
	_, ok := PostCategoryMap[id]
	return ok
}

type Post struct {
	ID         int    `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`
	CategoryID uint8  `gorm:"column:category_id;type:tinyint unsigned;not null;default:1;comment:Post category" json:"category_id"`
	UserID     int    `gorm:"column:user_id;type:int;not null;comment:User ID" json:"user_id"`
	Title      string `gorm:"column:title;type:varchar(100);not null;comment:Title" json:"title"`
	Summary    string `gorm:"column:summary;type:varchar(200);not null;comment:Post summary" json:"summary"`
	Slug       string `gorm:"column:slug;type:varchar(50);not null;comment:Post slug" json:"slug"`
	Content    string `gorm:"column:content;type:text;comment:Post content" json:"content"`
	Status     uint8  `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:Post status" json:"status"`

	Model
}
