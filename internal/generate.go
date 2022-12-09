package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"math/big"
	"math/rand"
	"unsafe"
)

func ByteArrayToInt(b []byte) int64 {
	return *(*int64)(unsafe.Pointer(&b[0]))
}

func randStringRunes(n uint8, seed int64) string {
	rand.Seed(seed)
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateFlag(key string, studentID string) string {
	// Generate flag based on the task key and studentID
	// Use key to calculate HMAC
	// If key size is not 64 bytes, SHA256 and padding used if less, just SHA256 if more
	// HMAC used as flag
	seedBytes, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	mac := hmac.New(sha256.New, seedBytes)
	mac.Write([]byte((studentID)))
	flag := mac.Sum(nil)
	return fmt.Sprintf("%x", flag)

}
func GenerateFlagSeed(taskSeed string, studentID string, length uint8) string {
	// Generate deterministic flag based on the task seed and studentID
	// Instead of HMAC, we need

	seedInt := new(big.Int)
	seedInt.SetString(taskSeed, 16)
	maxInt := new(big.Int)
	maxInt.SetString("7FFFFFFFFFFFFFFF", 16)
	if seedInt.Cmp(maxInt) == 1 {
		panic("Seed was too large (larger than int64)")
	}
	seedBytes, err := hex.DecodeString(taskSeed)
	if err != nil {
		panic(err)
	}
	mac := hmac.New(sha256.New, seedBytes)
	mac.Write([]byte((studentID)))
	secondSeed := mac.Sum(nil)[:15]
	byteToInt := ByteArrayToInt(secondSeed[:15])
	if byteToInt < 0 {
		panic("Seed for task generation was lower than zero, should not happen.")
	}
	flag := randStringRunes(length, byteToInt)
	return flag

}

func GenerateForSingleTask(cCtx *cli.Context, logger *zap.Logger) {
	task1Seed := "0ceae473f0c93e4c0dc5"
	//task1Seed := "ee90e1b33cd904fe8420"
	flag := GenerateFlag(task1Seed, cCtx.String("student"))
	//logger.Infof("ANSWER{%s}", flag)
	logger.Debug("Flag generated",
		zap.String("flag", flag),
	)
	//CreateCourseTasks(cCtx)
}
