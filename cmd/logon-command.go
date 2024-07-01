package cmd

import (
	"fmt"
	"os"

	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var logonCmd = &cobra.Command{
	Use:   "add-logon",
	Short: "Add a logon command to the Windows Sandbox configuration",
	Long:  "Add a logon command to the Windows Sandbox configuration",
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetString("name")
		c, err := winbox.Load(n + winbox.Ext)
		if err != nil {
			fmt.Println(err)
			return
		}

		l, _ := cmd.Flags().GetString("command")
		c.LogonCommand = winbox.Command{Command: l}

		f, err := os.Create(n + winbox.Ext)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := c.WriteXML(f); err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	logonCmd.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	logonCmd.Flags().StringP("command", "c", "", "Command to run on logon")

	rootCmd.AddCommand(logonCmd)
}
