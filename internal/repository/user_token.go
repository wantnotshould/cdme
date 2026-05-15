// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import (
	"context"
	"time"

	"code.cn/blog/internal/model"
	"github.com/google/uuid"
	"github.com/xiayoudi/ud"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
	*baseRepository
}

func NewUserTokenRepository(db *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{newBaseRepository(db)}
}

func (r *UserTokenRepository) WithTx(tx *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{newBaseRepository(tx)}
}

func (r *UserTokenRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&model.UserToken{})
}

func (r *UserTokenRepository) isValid(
	ctx context.Context,
	userID int,
	token string,
	tokenType uint8,
	jti uuid.UUID,
) (bool, error) {

	now := ud.Now()

	var exists bool

	err := r.baseQuery(ctx).
		Select("1").
		Where("jti = ?", jti[:]).
		Where("token = ?", token).
		Where("user_id = ?", userID).
		Where("token_type = ?", tokenType).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", now).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (r *UserTokenRepository) ValidateAccess(
	ctx context.Context,
	userID int,
	token string,
	jti uuid.UUID,
) (bool, error) {
	return r.isValid(ctx, userID, token, model.UserTokenTypeAccess, jti)
}

func (r *UserTokenRepository) ValidateRefresh(
	ctx context.Context,
	userID int,
	token string,
	jti uuid.UUID,
) (bool, error) {
	return r.isValid(ctx, userID, token, model.UserTokenTypeRefresh, jti)
}

func (r *UserTokenRepository) Add(ctx context.Context, userToken *model.UserToken) error {
	return r.baseQuery(ctx).Create(userToken).Error
}

func (r *UserTokenRepository) RevokeByJti(ctx context.Context, jti uuid.UUID) error {
	now := ud.Now()

	return r.baseQuery(ctx).
		Where("jti = ?", jti[:]).
		Where("revoked_at IS NULL").
		Update("revoked_at", now).Error
}

func (r *UserTokenRepository) RevokeBySessionID(
	ctx context.Context,
	userID int,
	sessionID uuid.UUID,
) error {

	now := ud.Now()

	return r.baseQuery(ctx).
		Where("user_id = ?", userID).
		Where("session_id = ?", sessionID[:]).
		Where("revoked_at IS NULL").
		Update("revoked_at", now).Error
}

func (r *UserTokenRepository) RevokeAllByUserID(ctx context.Context, userID int) error {
	now := ud.Now()

	return r.baseQuery(ctx).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Update("revoked_at", now).Error
}

func (r *UserTokenRepository) CleanupExpired(ctx context.Context) (int64, error) {
	return r.cleanupByCondition(ctx, "expires_at < ?", ud.Now().Add(-24*time.Hour))
}

func (r *UserTokenRepository) CleanupRevoked(ctx context.Context) (int64, error) {
	return r.cleanupByCondition(ctx, "revoked_at IS NOT NULL")
}

func (r *UserTokenRepository) cleanupByCondition(ctx context.Context, query string, args ...any) (int64, error) {
	var total int64

	for {
		result := r.db.WithContext(ctx).
			Model(&model.UserToken{}).
			Where(query, args...).
			Limit(500).
			Delete(&model.UserToken{})

		if result.Error != nil {
			return total, result.Error
		}

		affected := result.RowsAffected
		total += affected

		if affected < 500 {
			break
		}
	}

	return total, nil
}
