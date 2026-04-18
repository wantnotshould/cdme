// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "gorm.io/gorm"

type UserTokenRepository struct {
	*baseRepository
}

func NewUserTokenRepository(db *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{newBaseRepository(db)}
}

func (r *UserTokenRepository) WithTx(tx *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{newBaseRepository(tx)}
}
