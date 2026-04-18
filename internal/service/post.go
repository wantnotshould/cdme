// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package service

import (
	"context"
	"strings"

	"code.cn/blog/internal/dto/req"
	"code.cn/blog/internal/dto/resp"
	"code.cn/blog/internal/model"
	"code.cn/blog/internal/repository"
	"code.cn/blog/pkg/utils"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo}
}

func (s *PostService) WebList(
	ctx context.Context,
	param req.PostWebList,
) ([]resp.PostWebListItem, int64, error) {

	list, count, err := s.repo.WebList(ctx, param)
	if err != nil {
		return nil, 0, utils.Err("failed to retrieve post list")
	}

	posts := make([]resp.PostWebListItem, 0, len(list))

	for _, post := range list {
		posts = append(posts, resp.PostWebListItem{
			Title:     post.Title,
			Slug:      post.Slug,
			Summary:   post.Summary,
			CreatedAt: post.CreatedAt,
		})
	}

	return posts, count, nil
}

func (s *PostService) WebDetail(
	ctx context.Context,
	slug string,
) (*resp.PostWebDetail, error) {

	post, err := s.repo.WebDetail(ctx, slug)
	if err != nil {
		return nil, utils.Err("failed to retrieve post details")
	}

	if post == nil {
		return nil, utils.Err("post not found")
	}

	return &resp.PostWebDetail{
		Title:     post.Title,
		Slug:      post.Slug,
		Summary:   post.Summary,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}, nil
}

func validatePost(status uint8, categoryID uint8) error {

	if status > 0 && !model.IsValidPostStatus(status) {
		return utils.Err("invalid post status")
	}

	if categoryID > 0 && !model.IsValidPostCategory(categoryID) {
		return utils.Err("invalid post category")
	}

	return nil
}

func (s *PostService) Create(ctx context.Context, param req.PostCreate) error {
	if err := validatePost(param.Status, param.CategoryID); err != nil {
		return err
	}

	isDup, err := s.repo.Duplicate(ctx, param.UserID, param.Slug)
	if err != nil {
		return utils.Err("system busy")
	}

	if isDup {
		return utils.Err("slug already exists")
	}

	post := &model.Post{
		Title:      param.Title,
		CategoryID: model.PostCategory(param.CategoryID),
		Slug:       param.Slug,
		Summary:    param.Summary,
		Content:    param.Content,
		Status:     model.PostStatus(param.Status),
		UserID:     param.UserID,
	}

	return s.repo.Create(ctx, post)
}

func (s *PostService) Update(ctx context.Context, param req.PostUpdate) error {
	if err := validatePost(param.Status, param.CategoryID); err != nil {
		return err
	}

	info, err := s.repo.Info(ctx, req.InfoPathParams{
		InfoID: param.ID,
		UserID: param.UserID,
	})
	if err != nil {
		return utils.Err("system busy")
	}

	if info == nil {
		return utils.Err("post not found")
	}

	if !strings.EqualFold(info.Slug, param.Slug) {
		dup, err := s.repo.Duplicate(ctx, param.UserID, param.Slug, param.ID)
		if err != nil {
			return utils.Err("system busy")
		}
		if dup {
			return utils.Err("slug already exists")
		}
	}

	return s.repo.UpdateFields(ctx, param.ID, param.UserID, map[string]any{
		"title":       param.Title,
		"category_id": param.CategoryID,
		"summary":     param.Summary,
		"slug":        param.Slug,
		"content":     param.Content,
		"status":      param.Status,
	})
}
