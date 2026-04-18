// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "gorm.io/gorm"

type baseRepository struct {
	db *gorm.DB
}

func newBaseRepository(db *gorm.DB) *baseRepository {
	return &baseRepository{db}
}

func (b *baseRepository) ExecTx(fn func(tx *gorm.DB) error) error {
	return b.db.Transaction(fn)
}
