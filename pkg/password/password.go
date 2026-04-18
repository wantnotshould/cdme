// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"code.cn/blog/pkg/utils"
	"golang.org/x/crypto/argon2"
)

const (
	argonVariant   = "argon2id"
	argonVersion   = argon2.Version
	argonTime      = 3
	argonMemory    = 1 << 15
	argonThreads   = 2
	argonKeyLength = 32
	argonSaltLen   = 16
)

var hashSemaphore = make(chan struct{}, 8)

// Hash generates a hash string in PHC format
func Hash(password string) (string, error) {
	if password == "" {
		return "", utils.Err("Password cannot be empty")
	}

	// Limit concurrent hashing operations
	select {
	case hashSemaphore <- struct{}{}:
		defer func() { <-hashSemaphore }()
	case <-time.After(time.Second * 5):
		return "", utils.Err("Too many login requests, please try again later")
	}

	// Generate a random salt
	salt := make([]byte, argonSaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", utils.Wrap("Failed to generate salt", err)
	}

	// Generate the Argon2 hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLength,
	)

	// Encode salt and hash to Base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return the hash in PHC format
	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argonVariant, argonVersion, argonMemory, argonTime, argonThreads, b64Salt, b64Hash), nil
}

// Validate checks if the provided password matches the encoded hash
func Validate(password string, encodedHash string) (bool, error) {
	if password == "" || encodedHash == "" {
		return false, utils.Err("Both password and hash must be provided")
	}

	// Split the encoded hash into parts
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, utils.Err("Invalid hash format")
	}

	// Check the algorithm used
	if parts[1] != argonVariant {
		return false, utils.Err("Unsupported hashing algorithm")
	}

	// Check the version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, utils.Err("Failed to parse version")
	}
	if version != argonVersion {
		return false, utils.Err("Incompatible argon2 version")
	}

	// Parse parameters (memory, time, threads)
	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, utils.Err("Failed to parse parameters")
	}

	// Defensive check: avoid high memory DoS attacks
	if memory > 256*1024 {
		return false, utils.Err("This hash request uses too much memory")
	}

	// Enter the hashing section: get the semaphore
	hashSemaphore <- struct{}{}
	defer func() { <-hashSemaphore }()

	// Decode salt and hash from Base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, utils.Err("Failed to decode salt")
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, utils.Err("Failed to decode hash")
	}

	// Generate the comparison hash using the provided password and parameters
	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(decodedHash)),
	)

	// Compare the generated hash with the decoded hash using constant-time comparison
	return subtle.ConstantTimeCompare(comparisonHash, decodedHash) == 1, nil
}
