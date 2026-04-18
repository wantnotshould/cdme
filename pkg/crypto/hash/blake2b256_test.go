// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package hash

import "testing"

func TestHashBlake2b256(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"text", []byte("hello world")},
		{"binary", []byte{0x00, 0x01, 0x02}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := HashBlake2b256(tt.data)
			h2 := HashBlake2b256(tt.data)

			if len(h1) != 32 {
				t.Errorf("expected 32 bytes, got %d", len(h1))
			}

			if string(h1) != string(h2) {
				t.Error("hash should be deterministic")
			}
		})
	}
}

func TestHMACBlake2b256(t *testing.T) {
	key := []byte("secret-key")

	tests := [][]byte{
		[]byte("hello"),
		[]byte(""),
		[]byte("golang crypto test"),
	}

	for _, data := range tests {
		h1 := HMACBlake2b256(data, key)
		h2 := HMACBlake2b256(data, key)

		if string(h1) != string(h2) {
			t.Error("HMAC should be deterministic")
		}

		if len(h1) != 32 {
			t.Errorf("expected 32 bytes HMAC, got %d", len(h1))
		}
	}
}

func TestVerifyHMACBlake2b256(t *testing.T) {
	key := []byte("secret-key")
	data := []byte("hello world")

	sign := HMACBlake2b256Hex(data, key)

	if !VerifyHMACBlake2b256(data, key, sign) {
		t.Error("expected verification to pass")
	}

	// tampered case
	if VerifyHMACBlake2b256([]byte("tampered"), key, sign) {
		t.Error("expected verification to fail")
	}
}
