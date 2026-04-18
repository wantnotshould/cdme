// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package aes

import (
	"sync"
	"testing"
)

func resetForTest() {
	global = nil
	once = sync.Once{}
}

func TestAESGCM_AAD_Basic(t *testing.T) {
	resetForTest()

	Init([]byte("1234567890123456"))

	c := Global()

	plaintext := []byte("hello aad")

	aad1 := []byte("user:123")
	aad2 := []byte("user:456")

	enc, err := c.Encrypt(plaintext, aad1)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	dec, err := c.Decrypt(enc, aad1)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if string(dec) != string(plaintext) {
		t.Fatalf("mismatch: got %s want %s", dec, plaintext)
	}

	_, err = c.Decrypt(enc, aad2)
	if err == nil {
		t.Fatal("expected decrypt to fail with wrong AAD")
	}
}

func TestAESGCM_AAD_Tamper(t *testing.T) {
	resetForTest()

	Init([]byte("1234567890123456"))

	c := Global()

	plaintext := []byte("secure data")
	aad := []byte("tenant:A")

	enc, err := c.Encrypt(plaintext, aad)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	aadTampered := []byte("tenant:B")

	_, err = c.Decrypt(enc, aadTampered)
	if err == nil {
		t.Fatal("expected failure when AAD is tampered")
	}
}
