// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package redis

import (
	"context"
	"strconv"
	"time"
)

func (r *Redis) accessTokenKey(userID int) string {
	return r.key("access:" + strconv.Itoa(userID))
}

func (r *Redis) refreshTokenKey(userID int) string {
	return r.key("refresh:" + strconv.Itoa(userID))
}

func (r *Redis) SetAccessToken(
	ctx context.Context,
	userID int,
	token string,
	ttl time.Duration,
) error {

	return r.client.Set(
		ctx,
		r.accessTokenKey(userID),
		token,
		ttl,
	).Err()
}

func (r *Redis) GetAccessToken(
	ctx context.Context,
	userID int,
) (string, error) {

	return r.client.Get(
		ctx,
		r.accessTokenKey(userID),
	).Result()
}

func (r *Redis) SetRefreshToken(
	ctx context.Context,
	userID int,
	token string,
	ttl time.Duration,
) error {

	return r.client.Set(
		ctx,
		r.refreshTokenKey(userID),
		token,
		ttl,
	).Err()
}

func (r *Redis) GetRefreshToken(
	ctx context.Context,
	userID int,
) (string, error) {

	return r.client.Get(
		ctx,
		r.refreshTokenKey(userID),
	).Result()
}
