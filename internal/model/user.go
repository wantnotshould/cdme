// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package model

const (
	_ uint8 = iota
	UserStatusActive
	UserStatusDisabled
)

type User struct {
	ID           int    `gorm:"column:id;type:int;primaryKey;autoIncrement;comment:ID" json:"id"`
	Username     string `gorm:"column:username;type:varchar(26);not null;comment:Username" json:"username"`
	PasswordHash []byte `gorm:"column:password_hash;type:blob;not null;comment:Password hash" json:"-"`
	Status       uint8  `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:User status" json:"status"`

	Model
}
