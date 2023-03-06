package cmd

/*
Copyright © 2023 OUSPG ouspg@ouspg.org
*/


import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)
// Note that variables are defined package wide!

var defaultConfig = "course"
var defaultConfigPath = "."

// Loaded viper data
// Define here that can be accessed after parsing command-line parameters
var v *viper.Viper
var envPrefix = "GITHUB"
var logLevel ZapLogLevel
var enableJSON bool
var Logger *zap.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "graas",
	Short: "Grading Assistant",
	Long:  `Assistant for generating and reviewing student tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initConfig(cmd)
	},
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Map GitHub evironment variables to command names to keep access and modifiying easier with Cobra/
// Prefix is automatically removed from ENV variable (see var  envPrefix), suffix lowercased, and result compared to command argument names
// as result, following command names will match ENV variables, comparison is made with viper on initConfig
var (
    GitRepository = "repository" // "GITHUB_REPOSITORY"
    GitToken = "token" // "GITHUB_TOKEN"
    GitRef = "ref" // "GITHUB_REF"
    GitRefType = "ref-type" // "GITHUB_REF_TYPE"

)

func init() {
	// Global flags
    //
	cobra.OnInitialize(initLogger)
	rootCmd.PersistentFlags().StringP("config", "c", "", "Course config file (default is course.toml in current directory")
	rootCmd.PersistentFlags().StringP(GitRepository, "r", "", "Override GITHUB_REPOSITORY environment variable (target student)")
	rootCmd.MarkFlagRequired(GitRepository)
	rootCmd.PersistentFlags().String(GitToken, "", "Override GITHUB_TOKEN environment variable for GitHub authentication purposes")
	rootCmd.PersistentFlags().String(GitRefType, "", "Override GITHUB_REF_TYPE environment variable. The type of ref that triggered the GitHub Action workflow run.")
	rootCmd.PersistentFlags().String(GitRef, "", "Override GITHUB_REF environment variable. The fully-formed ref of the branch or tag that triggered the GitHub Actions workflow run.")
	rootCmd.PersistentFlags().VarP(&logLevel, "log", "l", "Set a log level. Available levels: debug, info, warn, error, dpanic, panic, fatal")
	rootCmd.PersistentFlags().BoolVarP(&enableJSON, "json", "j", false, "Enable JSON output")
	// Local flags
}

// ZapLogLevel overrides zapcore.Level package to support cobra variable flag directly
type ZapLogLevel struct {
	zapcore.Level
}

func (e *ZapLogLevel) Type() string {
	return "ZapLogLevel"
}

func initLogger() {
	// Init logger configuration for Zap
    // Currently only debug and production mode is supported with parameter logLevel
	var cfg zap.Config
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zapcore.LevelOf(logLevel))
	cfg.Level = atom
	if enableJSON {
		cfg.Encoding = "json"
	} else {
		cfg.Encoding = "console"
	}
	if cfg.Level.Level() == zap.DebugLevel {
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.Development = true
		cfg.DisableCaller = false
	} else {
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
		cfg.EncoderConfig.TimeKey = "" // disable timestamps
		cfg.Development = false
		cfg.DisableCaller = true
	}
	cfg.OutputPaths = []string{"stdout", "/tmp/logs"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	Logger = zap.Must(cfg.Build())
	defer Logger.Sync()
	Logger.Debug("Logger properly configured",
		zap.String("encoding", cfg.Encoding),
		zap.Bool("production", !cfg.Development),
	)
}

func initConfig(cmd *cobra.Command) error {
	var replaceHyphenWithCamelCase = true
	v = viper.New()
	v.SetConfigName(defaultConfig)
	v.AddConfigPath(defaultConfigPath)
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		Logger.Debug("Error when reading configuration file.",
			zap.String("error", err.Error()),
		)
		return err
	}
	cmd.Root().Flags().Lookup("")
    cmd.Flags().VisitAll(func(f *pflag.Flag) {
        // Determine the naming convention of the flags when represented in the config file
        configKey := f.Name
        // If using camelCase in the config file, replace hyphens with a camelCased string.
        // Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
        if replaceHyphenWithCamelCase {
            configKey = strings.ReplaceAll(f.Name, "-", "")
        }
        // Apply the viper config value to the flag when the flag is not set and viper has a value
        if !f.Changed && v.IsSet(configKey) {
            val := v.Get(configKey)
            err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
            if err != nil {
                Logger.Warn("Setting key in Viper failed.",
                zap.String("error", err.Error()))
            }
		}
	})
    return nil
}
