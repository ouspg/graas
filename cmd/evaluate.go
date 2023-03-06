package cmd

/*
Copyright © 2023 OUSPG ouspg@ouspg.org
*/

import (
	"fmt"

	"github.com/spf13/cobra"
)

// evaluateCmd represents the evaluate command
var evaluateCmd = &cobra.Command{
	Use:   "evaluate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("evaluate called")
	},
}

func init() {
	rootCmd.AddCommand(evaluateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// evaluateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// evaluateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
