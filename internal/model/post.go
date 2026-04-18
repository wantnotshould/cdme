// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package model

type PostStatus uint8
type PostCategory uint8

const (
	PostStatusDraft PostStatus = iota + 1
	PostStatusPublished
	PostStatusUnderReview
	PostStatusPrivate
)

const (
	PostCategorySharing PostCategory = iota + 1
	PostCategoryTech
	PostCategoryEssay
	PostCategoryTutorial
)

type Option struct {
	ID    uint8  `json:"id"`
	Label string `json:"label"`
}

var PostStatusOptions = []Option{
	{ID: uint8(PostStatusDraft), Label: "Draft"},
	{ID: uint8(PostStatusPublished), Label: "Published"},
	{ID: uint8(PostStatusUnderReview), Label: "Under Review"},
	{ID: uint8(PostStatusPrivate), Label: "Private"},
}

var PostCategoryOptions = []Option{
	{ID: uint8(PostCategorySharing), Label: "Sharing"},
	{ID: uint8(PostCategoryTech), Label: "Tech"},
	{ID: uint8(PostCategoryEssay), Label: "Essay"},
	{ID: uint8(PostCategoryTutorial), Label: "Tutorial"},
}

func toMap(opts []Option) map[uint8]Option {
	m := make(map[uint8]Option, len(opts))
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
	ID         int          `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`
	CategoryID PostCategory `gorm:"column:category_id;type:tinyint unsigned;not null;default:1;comment:Post category" json:"category_id"`
	UserID     int          `gorm:"column:user_id;type:int;not null;comment:User ID" json:"user_id"`
	Title      string       `gorm:"column:title;type:varchar(100);not null;comment:Title" json:"title"`
	Summary    string       `gorm:"column:summary;type:varchar(200);not null;comment:Post summary" json:"summary"`
	Slug       string       `gorm:"column:slug;type:varchar(50);not null;comment:Post slug" json:"slug"`
	Content    string       `gorm:"column:content;type:text;comment:Post content" json:"content"`
	Status     PostStatus   `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:Post status" json:"status"`

	Model
}
