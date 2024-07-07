package cmd

import (
	"fmt"
	"os"

	inputs "github.com/aboxofsox/winbox/tui/inputs"
	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var logonCmd = &cobra.Command{
	Use:   "add-logon",
	Short: "Add a logon command to the Windows Sandbox configuration",
	Long:  "Add a logon command to the Windows Sandbox configuration",
	Run: func(cmd *cobra.Command, args []string) {
		useTui, _ := cmd.Flags().GetBool("tui")
		if useTui {
			logonCommandWithTui()
		} else {
			logonCommandWithoutTui(cmd)
		}
	},
}

func logonCommandWithoutTui(cmd *cobra.Command) {
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
}

func logonCommandWithTui() {
	m := inputs.Show([]string{
		"Configuration Name",
		"Command",
	})

	tm := m.(inputs.Model)
	if tm.Inputs[0].Value() == "" {
		return
	}

	c, err := winbox.Load(tm.Inputs[0].Value() + winbox.Ext)
	if err != nil {
		panic(err)
	}

	if tm.Inputs[1].Value() == "" {
		fmt.Println("Command cannot be empty")
		return
	}
	c.AddLogonCommand(winbox.Command{
		Command: tm.Inputs[1].Value(),
	})
}

func init() {
	logonCmd.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	logonCmd.Flags().StringP("command", "c", "", "Command to run on logon")
	logonCmd.Flags().BoolP("tui", "u", false, "Use the TUI to add a logon command to the configuration")

	rootCmd.AddCommand(logonCmd)
}
