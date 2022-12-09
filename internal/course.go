package internal

import (
	"encoding/hex"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
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

func readCourseConfig(path string, logger *zap.Logger) (Course, toml.MetaData) {
	if _, err := os.Stat(path); err != nil {
		logger.Fatal("TOML configuration file does not exist.",
			zap.Error(err),
		)
	}
	var config Course
	meta, err := toml.DecodeFile(path, &config)
	if err != nil {
		logger.Fatal("Failed to decode configuration TOML.",
			zap.Error(err),
		)
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
	return config, meta

}

func getCourseInfo(path string, logger *zap.Logger) (Course, toml.MetaData) {
	config, meta := readCourseConfig(path, logger)
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
	return config, meta
}

func CreateCourseTasks(cCtx *cli.Context, logger *zap.Logger) {
	config, _ := getCourseInfo(cCtx.String("config"), logger)
	fmt.Println(config.Title)
}
