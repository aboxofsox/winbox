package cmd

import (
	clist "github.com/aboxofsox/winbox/tui/list"
	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a Windows Sandbox configuration",
	Long:  "Select a Windows Sandbox configuration",
	Run: func(cmd *cobra.Command, args []string) {
		name := clist.Show("Select a Windows Sandbox configuration", "Launching Windows Sandbox", clist.FindAllWSBFiles())
		if name == "" {
			return
		}

		err := winbox.Run(name)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
