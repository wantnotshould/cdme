// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "gorm.io/gorm"

type UserRepository struct {
	*baseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{newBaseRepository(db)}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{newBaseRepository(tx)}
}
