package cmd

import (
	"fmt"
	"os/exec"

	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Windows Sandbox configuration",
	Long:  "Run a Windows Sandbox configuration",
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetString("name")
		if n != "" {
			cmd := exec.Command("cmd", "/c", "start", n+winbox.Ext)
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return
			}
		} else {
			cmd := exec.Command("cmd", "/c", "start", winbox.WindowsSandboxPath)
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	runCmd.Flags().StringP("name", "N", "", "Name of the Windows Sandbox configuration")

	rootCmd.AddCommand(runCmd)
}
