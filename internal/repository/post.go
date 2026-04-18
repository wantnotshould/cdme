// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package repository

import (
	"context"
	"errors"

	"code.cn/blog/internal/dto/req"
	"code.cn/blog/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	*baseRepository
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{newBaseRepository(db)}
}

func (r *PostRepository) WithTx(tx *gorm.DB) *PostRepository {
	return &PostRepository{newBaseRepository(tx)}
}

func (r *PostRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&model.Post{})
}

func (r *PostRepository) buildQuery(db *gorm.DB, param req.PostList) *gorm.DB {
	db = db.Where("user_id = ?", param.UserID)

	if param.Keyword != "" {
		like := "%" + param.Keyword + "%"
		db = db.Where("(title LIKE ? OR summary LIKE ?)", like, like)
	}

	if param.Status > 0 {
		db = db.Where("status = ?", param.Status)
	}

	if param.CategoryID > 0 {
		db = db.Where("category_id = ?", param.CategoryID)
	}

	return db
}

func (r *PostRepository) executeList(
	query *gorm.DB,
	page, perPage int,
	selectFields []string,
) ([]model.Post, int64, error) {

	var (
		list  []model.Post
		count int64
	)

	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []model.Post{}, 0, nil
	}

	dataQuery := query.Order("created_at DESC")

	if len(selectFields) > 0 {
		dataQuery = dataQuery.Select(selectFields)
	} else {
		dataQuery = dataQuery.Select("id", "title", "slug", "summary", "created_at")
	}

	err := dataQuery.
		Scopes(paginateScope(page, perPage)).
		Find(&list).Error

	return list, count, err
}

func (r *PostRepository) List(ctx context.Context, param req.PostList) ([]model.Post, int64, error) {
	query := r.buildQuery(r.baseQuery(ctx), param)

	fields := []string{
		"id", "title", "slug", "category_id",
		"summary", "status", "created_at", "updated_at",
	}

	return r.executeList(query, param.Page, param.PerPage, fields)
}

func (r *PostRepository) Info(ctx context.Context, infoID, userID int) (*model.Post, error) {
	var info model.Post

	err := r.baseQuery(ctx).
		Where("id = ?", infoID).
		Where("user_id = ?", userID).
		Take(&info).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &info, nil
}

func (r *PostRepository) Duplicate(ctx context.Context, userID int, slug string, id ...int) (bool, error) {
	var exists bool

	q := r.baseQuery(ctx).
		Select("1").
		Where("user_id = ?", userID).
		Where("slug = ?", slug)

	if len(id) > 0 && id[0] > 0 {
		q = q.Where("id <> ?", id[0])
	}

	err := q.Limit(1).Scan(&exists).Error

	return exists, err
}

func (r *PostRepository) Create(ctx context.Context, p *model.Post) error {
	return r.baseQuery(ctx).Create(p).Error
}

func (r *PostRepository) UpdateFields(ctx context.Context, id, userID int, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	return r.baseQuery(ctx).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Updates(data).Error
}

func (r *PostRepository) BatchUpdateField(
	ctx context.Context,
	ids []int,
	userID int,
	column string,
	value any,
) error {

	if len(ids) == 0 {
		return nil
	}

	return r.baseQuery(ctx).
		Where("id IN ?", ids).
		Where("user_id = ?", userID).
		Update(column, value).Error
}

func (r *PostRepository) Delete(ctx context.Context, param req.DeletePathParams) error {
	return r.baseQuery(ctx).
		Where("id = ?", param.DeleteID).
		Where("user_id = ?", param.UserID).
		Delete(&model.Post{}).Error
}

func (r *PostRepository) BatchDelete(ctx context.Context, param req.BatchDelete) error {
	if len(param.DeleteIds) == 0 {
		return nil
	}

	return r.baseQuery(ctx).
		Where("id IN ?", param.DeleteIds).
		Where("user_id = ?", param.UserID).
		Delete(&model.Post{}).Error
}
