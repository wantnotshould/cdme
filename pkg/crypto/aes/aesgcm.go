// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"sync"
)

var (
	global *aesgcm
	once   sync.Once
)

type aesgcm struct {
	key []byte
	gcm cipher.AEAD
}

func Init(key []byte) {
	once.Do(func() {
		block, err := aes.NewCipher(key)
		if err != nil {
			log.Fatalf("aes init failed: new cipher error: %v", err)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Fatalf("aes init failed: gcm creation error: %v", err)
		}

		global = &aesgcm{
			key: key,
			gcm: gcm,
		}
	})
}

func Global() *aesgcm {
	if global == nil {
		panic("aes not initialized")
	}
	return global
}

func (a *aesgcm) Encrypt(plaintext []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, aad)

	result := make([]byte, 0, len(nonce)+len(ciphertext))
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func (a *aesgcm) Decrypt(data []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, aad)
}
