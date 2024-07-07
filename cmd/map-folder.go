package cmd

import (
	"fmt"
	"os"

	inputs "github.com/aboxofsox/winbox/tui/inputs"
	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var mapFolder = &cobra.Command{
	Use:   "map",
	Short: "Map a folder from the host to Windows Sandbox",
	Long:  "Map a folder from the host to Windows Sandbox",
	Run: func(cmd *cobra.Command, args []string) {
		useTui, _ := cmd.Flags().GetBool("tui")
		if useTui {
			mapFolderWithTui()
		} else {
			mapFolderWithoutTui(cmd)
		}
	},
}

func mapFolderWithoutTui(cmd *cobra.Command) {
	n, _ := cmd.Flags().GetString("name")
	h, _ := cmd.Flags().GetString("host")
	s, _ := cmd.Flags().GetString("sandbox")
	r, _ := cmd.Flags().GetBool("readonly")

	c, err := winbox.Load(n + winbox.Ext)
	if err != nil {
		fmt.Println(err)
		return
	}

	mf := winbox.MappedFolder{
		HostFolder:    h,
		SandboxFolder: s,
	}

	if r {
		mf.ReadOnly = true
	}

	c.MappedFolders = append(c.MappedFolders, mf)

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

func mapFolderWithTui() {
	m := inputs.Show([]string{
		"Configuration Name",
		"Host Folder",
		"Sandbox Folder",
		"Read-only (true/false)",
	})

	tm := m.(inputs.Model)
	if tm.Inputs[0].Value() == "" {
		return
	}

	c, err := winbox.Load(tm.Inputs[0].Value() + winbox.Ext)
	if err != nil {
		panic(err)
	}

	c.AddMappedFolder(winbox.MappedFolder{
		HostFolder:    tm.Inputs[1].Value(),
		SandboxFolder: tm.Inputs[2].Value(),
		ReadOnly:      isReadOnly(tm.Inputs[3].Value()),
	})
}

func isReadOnly(s string) bool {
	switch s {
	case "yes", "y", "true", "t":
		return true
	default:
		return false
	}
}

func init() {
	mapFolder.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	mapFolder.Flags().StringP("host", "H", "", "Host folder")
	mapFolder.Flags().StringP("sandbox", "S", "", "Sandbox folder")
	mapFolder.Flags().BoolP("readonly", "R", false, "Read-only")
	mapFolder.Flags().BoolP("tui", "u", false, "Use the TUI to map a folder")

	rootCmd.AddCommand(mapFolder)
}
