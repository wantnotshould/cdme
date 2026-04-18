// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import "gorm.io/gorm"

type PostRepository struct {
	*baseRepository
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{newBaseRepository(db)}
}

func (r *PostRepository) WithTx(tx *gorm.DB) *PostRepository {
	return &PostRepository{newBaseRepository(tx)}
}
