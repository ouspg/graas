/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "github.com/ouspg/graas/internal"
    "github.com/spf13/cobra"
    //    "github.com/spf13/pflag"
    "go.uber.org/zap"
    "math/big"
    "math/rand"
    "unsafe"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
//		fmt.Println(v.Get("title"))
        GenerateForSingleTask(cmd, Logger)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func ByteArrayToInt(b []byte) int64 {
	return *(*int64)(unsafe.Pointer(&b[0]))
}

func randStringRunes(n uint8, seed int64) string {
    // INSECURE deterministic RNG. The result can be reproduced with SEED
	rand.Seed(seed)
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateFlag(key string, studentID string) string {
	// Generate flag based on the task key and unique identifier
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
	// TODO Does not work well on very big seed at the moment

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

func GenerateForSingleTask(cmd *cobra.Command, logger *zap.Logger) {
	task1Seed := "0ceae473f0c93e4c0dc5"
	//task1Seed := "ee90e1b33cd904fe8420"
    // Repository is unique identifier for the student
    studentID, err :=  cmd.Flags().GetString(GitRepository)
    if err != nil || studentID == "" {

        logger.Fatal("Unique student identifier required for task generation",
            zap.String("identifier", studentID),
            )
    }
    flag := GenerateFlag(task1Seed, studentID)
	//logger.Infof("ANSWER{%s}", flag)
	logger.Debug("Flag generated",
		zap.String("flag", flag),
	)
	internal.CreateCourseTasks(v, Logger)
}
