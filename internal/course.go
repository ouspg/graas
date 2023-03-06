package internal

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"path/filepath"
)

type (
	Course struct {
		Title    string
		TasksDir string
		Tasks    map[string][]struct {
			Number   uint16
			Key      string
			FlagType string
		}
	}
)

func readCourseConfig(viperData *viper.Viper, logger *zap.Logger) *Course {

	config := &Course{}
	err := viperData.Unmarshal(config)
	if err != nil {
		logger.Fatal("unable to decode into config struct",
			zap.String("error", err.Error()))
	}

	if config.TasksDir == "" {
		logger.Fatal("Task directory path must be provided.")
	}

	if !filepath.IsAbs(config.TasksDir) {
		logger.Fatal("Task directory path is not absolute.")
	}

	// Iterate over all weeks and tasks in the course
	// No duplicate keys allowed
	// Check key format and mandatory fields
	var uniqueKeys = make(map[string]bool)
	for week, value := range config.Tasks {
		for _, task := range value {
			if task.FlagType == "" {
				logger.Fatal("Flag type must be set.",
					zap.String("week", week),
					zap.Uint16("taskNumber", task.Number),
				)
			}
			if task.Number == 0 {
				logger.Fatal("Task number must be set. Other than zero.",
					zap.String("week", week),
					zap.Uint16("taskNumber", task.Number),
				)
			}
			if task.Key == "" || len(task.Key) < 64 {
				logger.Fatal("Key must be set and be at least 64 characters long hexadecimal string.",
					zap.String("week", week),
					zap.Uint16("taskNumber", task.Number),
				)
			}
			dst := make([]byte, hex.DecodedLen(len(task.Key)))
			if _, err := hex.Decode(dst, []byte(task.Key)); err != nil {
				// NOTE error can leak information about partial key
				logger.Fatal("Key must be valid hexadecimal string.",
					zap.String("week", week),
					zap.Uint16("taskNumber", task.Number),
					zap.Error(err),
				)
			}
			if uniqueKeys[task.Key] {
				logger.Fatal("All keys are not unique for the tasks. Provided occurrence has repeated key.",
					// Keys are not printed to avoid leaking information
					zap.String("week", week),
					zap.Uint16("taskNumber", task.Number),
				)
			}
			uniqueKeys[task.Key] = true
		}
	}
	return config

}

func getCourseInfo(v *viper.Viper, logger *zap.Logger) *Course {
	config := readCourseConfig(v, logger)
	var (
		numTasks = 0
	)
	for _, value := range config.Tasks {
		numTasks += len(value)
	}
	logger.Info("Course data loaded successfully.")
	logger.Debug("Course information:",
		zap.String("title", config.Title),
		zap.String("taskDirectory", config.TasksDir),
		zap.Int("numberWeeks", len(config.Tasks)),
		zap.Int("numberTotalTasks", numTasks),
	)
	return config
}

func CreateCourseTasks(cmd *cobra.Command, v *viper.Viper, logger *zap.Logger) {
	config := getCourseInfo(v, logger)
	fmt.Println(config.Title)
}
func CreateSingleTask(v *viper.Viper, student string, weekNro int8, taskNro int8, logger *zap.Logger) {
	// Create full task, include build phase. Currently just creates a flag without any specific format
	// TODO add support for Flag types

	config := getCourseInfo(v, logger)
	fmt.Println(config.Title)
	taskKey := config.Tasks[fmt.Sprintf("week%d", weekNro)][taskNro-1].Key
	flag := GenerateSecureFlag(taskKey, student, logger)
	// TODO remove printing anything related to flag or key or disallow DEBUG mode on CI by configuring GitHub runner
	// Otherwise might leak information
	logger.Debug("Flag generated",
		zap.String("flag", flag),
	)

}
