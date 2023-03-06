package cmd

/*
Copyright © 2023 OUSPG ouspg@ouspg.org
*/

import (
	"github.com/ouspg/graas/internal"
	"github.com/spf13/cobra"
	//    "github.com/spf13/pflag"
	"go.uber.org/zap"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		//		fmt.Println(v.Get("title"))
		weekNro, err := cmd.Flags().GetInt8("week")
		if err != nil {
			return err
		}
		taskNro, err := cmd.Flags().GetInt8("task")
		if err != nil {
			return err
		}
		if weekNro > 0 && taskNro > 0 {
			GenerateForSingleTask(cmd, weekNro, taskNro, Logger)
		}
		// TODO generate full week, all weeks?
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Maybe create custom struct for flag types
	generateCmd.LocalFlags().StringP("type", "t", "env", "Set flag type to generate")
	generateCmd.LocalFlags().BoolP("flag", "f", false, "Generate flak only (not full task)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GenerateForSingleTask(cmd *cobra.Command, weekNro int8, taskNro int8, logger *zap.Logger) {
	// Repository is unique identifier for the student
	student, err := cmd.Flags().GetString(GitRepository)
	if err != nil || student == "" {

		logger.Fatal("Unique student identifier required for task generation",
			zap.String("identifier", student),
		)
	}
	// Priority of configurations: cmd > viper > defaults
	internal.CreateSingleTask(v, student, weekNro, taskNro, logger)
}
