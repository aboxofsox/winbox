package cmd

import (
	clist "github.com/aboxofsox/winbox/tui/list"
	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a Windows Sandbox configuration",
	Long:  "Select a Windows Sandbox configuration",
	Run: func(cmd *cobra.Command, args []string) {
		clist.Show()
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
