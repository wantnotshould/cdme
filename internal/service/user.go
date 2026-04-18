// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package service

import (
	"context"

	"code.cn/blog/conf"
	"code.cn/blog/internal/auth/token"
	"code.cn/blog/internal/cache/redis"
	"code.cn/blog/internal/consts"
	"code.cn/blog/internal/dto/req"
	"code.cn/blog/internal/dto/resp"
	"code.cn/blog/internal/model"
	"code.cn/blog/internal/repository"
	"code.cn/blog/pkg/crypto/hash"
	"code.cn/blog/pkg/password"
	"code.cn/blog/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	repo          *repository.UserRepository
	userTokenRepo *repository.UserTokenRepository
}

func NewUserService(
	repo *repository.UserRepository,
	userTokenRepo *repository.UserTokenRepository,
) *UserService {
	return &UserService{repo, userTokenRepo}
}

func (s *UserService) syncRedis(ctx context.Context, userID int, atHash, rtHash string) {
	redis.DB().SetAccessToken(ctx, userID, atHash, consts.ATDuration)
	redis.DB().SetRefreshToken(ctx, userID, rtHash, consts.RTDuration)
}

func (s *UserService) issueTokens(
	ctx context.Context,
	user *model.User,
	ip, ua string,
	oldSessionID uuid.UUID,
) (*token.Response, error) {
	sessionID := uuid.New()
	accessJti := uuid.New()
	refreshJti := uuid.New()

	param := token.Param{
		UserID:    user.ID,
		IP:        ip,
		UserAgent: ua,
		SessionID: sessionID,
	}

	res, err := token.Generate(param, accessJti, refreshJti)
	if err != nil {
		return nil, utils.Err("failed to generate tokens")
	}

	atHash := hash.HMACBlake2b256Hex([]byte(res.AccessToken), []byte(conf.Get().Hash.Key))
	rtHash := hash.HMACBlake2b256Hex([]byte(res.RefreshToken), []byte(conf.Get().Hash.Key))

	err = s.userTokenRepo.ExecTx(func(tx *gorm.DB) error {
		repo := s.userTokenRepo.WithTx(tx)

		// revoke logic
		if oldSessionID == uuid.Nil {
			if err = repo.RevokeAllByUserID(ctx, user.ID); err != nil {
				return err
			}
		} else {
			if err = repo.RevokeBySessionID(ctx, user.ID, oldSessionID); err != nil {
				return err
			}
		}

		base := model.UserToken{
			UserID:    user.ID,
			SessionID: sessionID[:],
			IP:        ip,
			UserAgent: ua,
		}

		at := base
		at.TokenType = model.UserTokenTypeAccess
		at.Token = atHash
		at.Jti = accessJti[:]
		at.ExpiresAt = res.AccessExpiresAt

		if err = repo.Add(ctx, &at); err != nil {
			return err
		}

		rt := base
		rt.TokenType = model.UserTokenTypeRefresh
		rt.Token = rtHash
		rt.Jti = refreshJti[:]
		rt.ExpiresAt = res.RefreshExpiresAt

		return repo.Add(ctx, &rt)
	})

	if err != nil {
		return nil, utils.Err("failed to persist tokens")
	}

	s.syncRedis(ctx, user.ID, atHash, rtHash)

	return res, nil
}

func (s *UserService) Login(ctx context.Context, param req.UserLogin) (*token.Response, error) {
	info, err := s.repo.InfoByUsername(ctx, param.Username)
	if err != nil || info == nil {
		return nil, utils.Err("invalid username or password")
	}

	if info.Status == model.UserStatusDisabled {
		return nil, utils.Err("user disabled")
	}

	ok, err := password.Validate(param.Password, string(info.PasswordHash))
	if err != nil || !ok {
		return nil, utils.Err("invalid username or password")
	}

	return s.issueTokens(ctx, info, param.IP, param.UserAgent, uuid.Nil)
}

func (s *UserService) Profile(ctx context.Context, userID int) (*resp.UserProfile, error) {
	if userID <= 0 {
		return nil, utils.Err("invalid user id")
	}

	info, err := s.repo.InfoByID(ctx, userID)
	if err != nil {
		return nil, utils.Err("failed to get user profile")
	}

	if info == nil {
		return nil, utils.Err("user not found")
	}

	return &resp.UserProfile{
		Username: info.Username,
	}, nil
}

func (s *UserService) RefreshToken(
	ctx context.Context,
	param req.UserRefreshToken,
	claims *token.Claims,
) (*token.Response, error) {

	userID := claims.DecryptedPayload.UserID

	storedHash, err := redis.DB().GetRefreshToken(ctx, userID)

	if err != nil {
		return nil, utils.Err("system busy")
	}

	if storedHash == "" || storedHash != param.RefreshToken {
		return nil, utils.Err("session expired")
	}

	valid, err := s.userTokenRepo.ValidateRefresh(
		ctx,
		userID,
		param.RefreshToken,
		claims.DecryptedPayload.Jti,
	)
	if err != nil || !valid {
		return nil, utils.Err("session invalid")
	}

	info, err := s.repo.InfoByID(ctx, claims.DecryptedPayload.UserID)
	if err != nil || info == nil {
		return nil, utils.Err("user invalid")
	}

	if info.Status == model.UserStatusDisabled {
		return nil, utils.Err("user disabled")
	}

	return s.issueTokens(ctx, info, param.IP, param.UserAgent, claims.SessionID)
}

func (s *UserService) Logout(ctx context.Context, userID int, sessionID uuid.UUID) error {
	if err := s.userTokenRepo.RevokeBySessionID(ctx, userID, sessionID); err != nil {
		return utils.Err("logout failed")
	}
	return nil
}
