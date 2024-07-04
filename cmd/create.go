package cmd

import (
	"fmt"
	"os"

	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a new Windows Sandbox configuration",
	Long:  "Create a new Windows Sandbox configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		c := &winbox.Configuration{}
		c.VGpu, _ = cmd.Flags().GetString("vGpu")
		c.Networking, _ = cmd.Flags().GetString("networking")
		c.AudioInput, _ = cmd.Flags().GetString("audio")
		c.VideoInput, _ = cmd.Flags().GetString("video")
		c.ProtectedClient, _ = cmd.Flags().GetString("protected")
		c.PrinterRedirection, _ = cmd.Flags().GetString("printer")
		c.ClipboardRedirection, _ = cmd.Flags().GetString("clipboard")
		c.MemoryInMB, _ = cmd.Flags().GetInt("memory")

		if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
			config, err := winbox.LoadWinboxConfig()
			if err != nil {
				fmt.Println(err)
				return
			}
			winbox.WindowsSandboxPath = config.WindowsSandboxPath
		}
		name, _ := cmd.Flags().GetString("name")
		f, err := os.Create(name + winbox.Ext)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		if err := c.WriteXML(f); err != nil {
			fmt.Println(err)
			return
		}

	},
}

func init() {
	create.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	create.Flags().StringP("vGpu", "v", "Disable", "Enable or disable vGPU")
	create.Flags().StringP("networking", "n", "Default", "Networking configuration")
	create.Flags().StringP("audio", "a", "Disable", "Audio input")
	create.Flags().StringP("video", "i", "Disable", "Video input")
	create.Flags().StringP("protected", "p", "Disable", "Protected client")
	create.Flags().StringP("printer", "r", "Disable", "Printer redirection")
	create.Flags().StringP("clipboard", "c", "Disable", "Clipboard redirection")
	create.Flags().IntP("memory", "m", 8*1024, "Memory in MB")

	rootCmd.AddCommand(create)
}
