package tui

import (
	"encoding/json"
	"fmt"
	"github.com/ouspg/grass/internal"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"os"
	"sort"
)

var version = "0.1"

func configLogger() zap.Config {
	//"initialFields": {"version": "%s"},
	rawJSON := []byte(`{
	  "level": "debug",
      "development": true,
      "disableCaller": false,
	  "encoding": "console",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)
	//formJSON := fmt.Sprintf(string(rawJSON), version)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func Run() {
	// For Posix short combination, use UseShortOptionHandler https://github.com/urfave/cli/pull/684
	app := &cli.App{
		UseShortOptionHandling: true,
		Name:                   "grass",
		Usage:                  "Grading Assistant",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`, contains course configuration",
			},
			&cli.StringFlag{
				Name:    "student",
				Aliases: []string{"s"},
				Usage:   "Unique identifier for the student, used as seed for generating and reviewing tasks.",
				EnvVars: []string{"GITHUB_REPOSITORY"},
			},
			&cli.StringFlag{
				Name: "gauth",
				//Aliases: []string{"s"},
				Usage:   "Default GitHub authentication token",
				EnvVars: []string{"GITHUB_TOKEN"},
			},
			&cli.StringFlag{
				Name: "ref-type",
				//Aliases: []string{"s"},
				Usage:   "The type of ref that triggered the GitHub Action workflow run.",
				EnvVars: []string{"GITHUB_REF_TYPE"},
			},
			&cli.StringFlag{
				Name: "ref",
				//Aliases: []string{"s"},
				Usage:   "The fully-formed ref of the branch or tag that triggered the GitHub Actions workflow run.",
				EnvVars: []string{"GITHUB_REF"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "evaluate",
				//Aliases: []string{"c"},
				Usage: "Evaluate task(s)",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "task",
						Aliases: []string{"t"},
						Usage:   "Specify task number.",
					},
					&cli.StringFlag{
						Name:  "type",
						Value: "markdown",
						Usage: "Target format for evaluation",
					},
				},
				Action: func(*cli.Context) error {
					fmt.Print("Voila")
					return nil
				},
			},
			{
				Name: "generate",
				//Aliases: []string{"a"},
				Usage: "Generate task(s)",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "week",
						Aliases: []string{"w"},
						Usage:   "Specify week.",
					},
					&cli.UintFlag{
						Name:    "task",
						Aliases: []string{"t"},
						Usage:   "Specify task number.",
					},
					&cli.BoolFlag{
						Name:    "flag",
						Aliases: []string{"f"},
						Usage:   "Generates full flag.",
					},
					&cli.BoolFlag{
						Name:  "type",
						Usage: "Type of the flag, e.g. export ENV, create .txt file or binary",
					},
				},
				Action: func(cCtx *cli.Context) error {
					logger := zap.Must(configLogger().Build())
					defer logger.Sync()
					logger.Debug("Evaluated parameters..",
						// Structured context as strongly typed Field values.
						zap.String("student", cCtx.String("student")),
						zap.String("config", cCtx.String("config")),
					)
					if cCtx.Uint("task") != 0 {
						//internal.GenerateForSingleTask(cCtx, logger)
						internal.CreateCourseTasks(cCtx, logger)
					}
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
