package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aboxofsox/eval"
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

	items := make([]string, len(c.Keys()))
	for i, k := range c.Keys() {
		items[i] = format(k)
	}
	tm := inputs.Show(items)
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

	exp := m.Inputs[3].Value()
	mem, err := eval.Eval(exp)
	if err != nil {
		fmt.Println(err)
		return

	}
	c.MemoryInMB = mem

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

func format(s string) string {
	var res string
	for i, c := range s {
		if i == 0 {
			res += strings.ToUpper(string(c))
			continue
		}
		if isCapital(c) && !isCapital(rune(s[i-1])) {
			res += " " + string(c)
		} else {
			res += string(c)
		}
	}
	return res
}

func indexOf(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

func isCapital(s rune) bool {
	return s >= 'A' && s <= 'Z'
}

func init() {
	create.Flags().BoolP("tui", "u", false, "Use the TUI to create a configuration")
	create.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	create.Flags().StringP("vGpu", "g", "Disable", "Enable or disable vGPU")
	create.Flags().StringP("networking", "n", "Default", "Networking configuration")
	create.Flags().StringP("audio", "a", "Disable", "Audio input")
	create.Flags().StringP("video", "v", "Disable", "Video input")
	create.Flags().StringP("protected", "p", "Disable", "Protected client")
	create.Flags().StringP("printer", "r", "Disable", "Printer redirection")
	create.Flags().StringP("clipboard", "c", "Disable", "Clipboard redirection")
	create.Flags().IntP("memory", "m", 8*1024, "Memory in MB")

	rootCmd.AddCommand(create)
}
