// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package hash

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/blake2b"
)

func HashBlake2b256(data []byte) []byte {
	h, _ := blake2b.New256(nil)
	h.Write(data)
	return h.Sum(nil)
}

func HashBlake2b256Hex(data []byte) string {
	return hex.EncodeToString(HashBlake2b256(data))
}

func HMACBlake2b256(data []byte, key []byte) []byte {
	mac := hmac.New(func() hash.Hash {
		h, _ := blake2b.New256(nil)
		return h
	}, key)

	mac.Write(data)
	return mac.Sum(nil)
}

func HMACBlake2b256Hex(data []byte, key []byte) string {
	return hex.EncodeToString(HMACBlake2b256(data, key))
}

func VerifyHMACBlake2b256(data []byte, key []byte, expectedHex string) bool {
	expected, err := hex.DecodeString(expectedHex)
	if err != nil {
		return false
	}

	actual := HMACBlake2b256(data, key)
	return hmac.Equal(actual, expected)
}
