package cmd

import (
	"fmt"
	"os"

	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var mapFolder = &cobra.Command{
	Use:   "map",
	Short: "Map a folder from the host to Windows Sandbox",
	Long:  "Map a folder from the host to Windows Sandbox",
	Run: func(cmd *cobra.Command, args []string) {
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

	},
}

func init() {
	mapFolder.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	mapFolder.Flags().StringP("host", "H", "", "Host folder")
	mapFolder.Flags().StringP("sandbox", "S", "", "Sandbox folder")
	mapFolder.Flags().BoolP("readonly", "R", false, "Read-only")

	rootCmd.AddCommand(mapFolder)
}
