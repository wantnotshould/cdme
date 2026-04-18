// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import (
	"context"
	"errors"

	"code.cn/blog/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	*baseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{newBaseRepository(db)}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{newBaseRepository(tx)}
}

func (r *UserRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&model.User{})
}

func (r *UserRepository) InfoByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	err := r.baseQuery(ctx).
		Where("id = ?", id).
		Take(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) InfoByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.baseQuery(ctx).
		Where("username = ?", username).
		Take(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.baseQuery(ctx).Create(user).Error
}