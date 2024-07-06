package cmd

import (
	"fmt"
	"os"
	"strconv"

	inputs "github.com/aboxofsox/winbox/tui/inputs"
	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var create = &cobra.Command{
	Use:   "create",
	Short: "Create a new Windows Sandbox Configuration",
	Long:  "Create a new Windows Sandbox Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		useTui, _ := cmd.Flags().GetBool("tui")
		if useTui {
			createWithTui()
		} else {
			createWithoutTui(cmd)
		}
	},
}

func createWithoutTui(cmd *cobra.Command) {
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
}

func createWithTui() {
	c := &winbox.Configuration{}

	tm := inputs.Show(c.Keys())
	m, ok := tm.(inputs.Model)
	if !ok {
		fmt.Println("oh no")
		return
	}

	n := m.Inputs[0].Value()
	if n == "" {
		return
	}
	c.AudioInput = m.Inputs[1].Value()
	c.ClipboardRedirection = m.Inputs[2].Value()

	// it would be nice to handle equations
	// i.e. 8 * 1024 or 1024 + 1024 + 512
	// the Shunting Yard algorithm comes to mind
	mimb, err := strconv.Atoi(m.Inputs[3].Value())
	if err == nil {
		c.MemoryInMB = mimb
	}

	c.Networking = m.Inputs[4].Value()
	c.PrinterRedirection = m.Inputs[5].Value()
	c.ProtectedClient = m.Inputs[6].Value()
	c.VGpu = m.Inputs[7].Value()
	c.VideoInput = m.Inputs[8].Value()

	f, err := os.Create(n + winbox.Ext)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = c.WriteXML(f)
	if err != nil {
		panic(err)
	}
}

func init() {
	create.Flags().BoolP("tui", "u", false, "Use the TUI to create a configuration")
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
