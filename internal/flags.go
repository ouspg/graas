package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"unsafe"
)

func ByteArrayToInt(b []byte) int64 {
	return *(*int64)(unsafe.Pointer(&b[0]))
}

func randStringRunes(n uint8, seed int64) string {
	// INSECURE deterministic RNG. The result can be reproduced with SEED
	// Just sample for now
	rand.Seed(seed)
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateSecureFlag(key string, student string, logger *zap.Logger) string {
	// Generate flag based on the task key and unique identifier (student repository name?)
	// Use key to calculate HMAC,
	// HMAC should be secure enough for known-plaintext attacks (advisor knows repository name)
	// If key size is not 64 bytes, SHA256 and padding used if less, just SHA256 if more
	// HMAC used as flag
	seedBytes, err := hex.DecodeString(key)
	if err != nil {
		logger.Panic("Failed to create Key from hex string",
			zap.String("error", err.Error()))
	}
	mac := hmac.New(sha256.New, seedBytes)
	mac.Write([]byte((student)))
	flag := mac.Sum(nil)
	return fmt.Sprintf("%x", flag)

}
