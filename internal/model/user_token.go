// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package model

import "time"

const (
	_ uint8 = iota
	UserTokenTypeAccess
	UserTokenTypeRefresh
)

type UserToken struct {
	ID        int        `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`
	TokenType uint8      `gorm:"column:token_type;type:tinyint unsigned;not null;default:1;comment:Token type" json:"token_type"`
	UserID    int        `gorm:"column:user_id;type:int;not null;comment:User ID" json:"user_id"`
	Token     string     `gorm:"column:token;type:char(64);not null;comment:Token hash" json:"token"`
	SessionID []byte     `gorm:"column:session_id;type:binary(16);not null;comment:Session ID" json:"session_id"`
	Jti       []byte     `gorm:"column:jti;type:binary(16);not null;comment:JWT unique identifier" json:"jti"`
	ExpiresAt time.Time  `gorm:"column:expires_at;type:datetime(3);not null;comment:Expiration time" json:"expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at;comment:Revocation time" json:"revoked_at,omitempty"`
	IP        string     `gorm:"column:ip;type:varchar(45);comment:Login IP (supports IPv6)" json:"ip"`
	UserAgent string     `gorm:"column:user_agent;type:varchar(512);comment:Browser/device information" json:"user_agent"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime(3);comment:Creation time" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime(3);comment:Update time" json:"updated_at"`
}
