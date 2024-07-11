package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.4.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of winbox",
	Long:  "Print the version of winbox",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)

	},
}

func init() {
	versionCmd.Flags().BoolP("semver", "s", false, "Print semver version")
	rootCmd.AddCommand(versionCmd)
}
