// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package token

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"code.cn/blog/conf"
	"code.cn/blog/internal/consts"
	"code.cn/blog/pkg/crypto/aes"
	"code.cn/blog/pkg/crypto/hash"
	"code.cn/blog/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	errInvalidToken = utils.Err("invalid token")
	errTokenExpired = utils.Err("token expired")
)

type Plaintext struct {
	UserID     int       `json:"user_id"`
	Jti        uuid.UUID `json:"jti"`
	ClientHash string    `json:"client_hash"`
}

type Claims struct {
	jwt.RegisteredClaims
	EncryptedPayload string    `json:"encrypted_payload"`
	SessionID        uuid.UUID `json:"session_id"`
	DecryptedPayload Plaintext `json:"-"`
}

type Param struct {
	UserID    int       `json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	SessionID uuid.UUID `json:"session_id"`
}

type Response struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func encrypt(p Plaintext) (string, error) {
	plainBytes, err := json.Marshal(p)
	if err != nil {
		return "", utils.Wrap("encrypt: marshal failed", err)
	}

	cipherBytes, err := aes.Global().Encrypt(
		plainBytes,
		[]byte(conf.Get().AESGCM.AAD),
	)
	if err != nil {
		return "", utils.Wrap("encrypt: aes encrypt failed", err)
	}

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func decrypt[T any](cipherText string) (*T, error) {
	rawBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, utils.Wrap("decrypt: base64 decode failed", err)
	}

	plainBytes, err := aes.Global().Decrypt(
		rawBytes,
		[]byte(conf.Get().AESGCM.AAD),
	)
	if err != nil {
		return nil, utils.Wrap("decrypt: aes decrypt failed", err)
	}

	var out T
	if err := json.Unmarshal(plainBytes, &out); err != nil {
		return nil, utils.Wrap("decrypt: json unmarshal failed", err)
	}

	return &out, nil
}

func Generate(param Param, accessJti, refreshJti uuid.UUID) (*Response, error) {
	now := utils.Now()

	accessExp := now.Add(consts.ATDuration)
	refreshExp := now.Add(consts.RTDuration)

	clientHash := hash.HMACBlake2b256Hex(
		[]byte(param.IP+param.UserAgent),
		[]byte(conf.Get().Hash.Key),
	)

	// Access Token Payload
	accessPayloadPlain := Plaintext{
		UserID:     param.UserID,
		Jti:        accessJti,
		ClientHash: clientHash,
	}

	accessPayloadEncrypted, err := encrypt(accessPayloadPlain)
	if err != nil {
		return nil, utils.Err("failed to encrypt access token")
	}

	accessClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", param.UserID),
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessJti.String(),
		},
		SessionID:        param.SessionID,
		EncryptedPayload: accessPayloadEncrypted,
	}

	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	).SignedString([]byte(conf.Get().JWT.Key))

	if err != nil {
		return nil, err
	}

	// Refresh Token Payload
	refreshPayloadPlain := Plaintext{
		UserID:     param.UserID,
		Jti:        refreshJti,
		ClientHash: clientHash,
	}

	refreshPayloadEncrypted, err := encrypt(refreshPayloadPlain)
	if err != nil {
		return nil, utils.Err("failed to encrypt refresh token")
	}

	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", param.UserID),
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshJti.String(),
		},
		SessionID:        param.SessionID,
		EncryptedPayload: refreshPayloadEncrypted,
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString([]byte(conf.Get().JWT.Key))

	if err != nil {
		return nil, err
	}

	return &Response{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func Parse(tokenString, ip, userAgent string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.Get().JWT.Key), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errTokenExpired
		}
		return nil, errInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errInvalidToken
	}

	if claims.EncryptedPayload != "" {
		decryptedPayload, err := decrypt[Plaintext](claims.EncryptedPayload)
		if err != nil {
			return nil, errInvalidToken
		}

		if fmt.Sprintf("%d", decryptedPayload.UserID) != claims.Subject {
			return nil, utils.Err("token identity mismatch")
		}

		claims.DecryptedPayload = *decryptedPayload

		valid := hash.VerifyHMACBlake2b256(
			[]byte(ip+userAgent),
			[]byte(conf.Get().Hash.Key),
			decryptedPayload.ClientHash,
		)

		if !valid {
			return nil, utils.Err("client binding mismatch")
		}
	}

	return claims, nil
}
