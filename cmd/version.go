package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	versionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Border(lipgloss.RoundedBorder()).Padding(1, 2).Margin(1, 2)
)

var version = "0.4.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of winbox",
	Long:  "Print the version of winbox",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionStyle.Render(version))

	},
}

func init() {
	versionCmd.Flags().BoolP("semver", "s", false, "Print semver version")
	rootCmd.AddCommand(versionCmd)
}
