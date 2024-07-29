package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	inputs "github.com/aboxofsox/winbox/tui/inputs"
	clist "github.com/aboxofsox/winbox/tui/list"
	"github.com/aboxofsox/winbox/winbox"
	"github.com/spf13/cobra"
)

var (
	SandboxAccountPath = "C:\\Users\\WDAGUtilityAccount"
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
	rm, _ := cmd.Flags().GetBool("remove")

	if rm {
		removeMappedFolderWithTui(n)
		return
	}

	c, err := winbox.Load(n + winbox.Ext)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.Contains(h, "$env:") {
		h = replaceEnv(h)
	}
	if strings.Contains(h, "%") {
		h = replaceOldEnv(h)
	}
	if strings.Contains(s, "$env:") {
		s = replaceEnv(s)
	}
	if strings.Contains(s, "%") {
		s = replaceOldEnv(s)
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

	hostFolder := tm.Inputs[1].Value()
	sandboxFolder := tm.Inputs[2].Value()

	if strings.Contains(hostFolder, "$env:") {
		hostFolder = replaceEnv(hostFolder)
	}
	if strings.Contains(hostFolder, "%") {
		hostFolder = replaceOldEnv(hostFolder)
	}

	if strings.Contains(sandboxFolder, "$SandboxUser") {
		sandboxFolder = strings.Replace(sandboxFolder, "$SandboxUser", SandboxAccountPath, 1)
	}
	if strings.Contains(sandboxFolder, "$env:") {
		sandboxFolder = replaceEnv(sandboxFolder)
	}

	c.AddMappedFolder(winbox.MappedFolder{
		HostFolder:    hostFolder,
		SandboxFolder: sandboxFolder,
		ReadOnly:      isReadOnly(tm.Inputs[3].Value()),
	})

	f, err := os.Create(tm.Inputs[0].Value() + winbox.Ext)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := c.WriteXML(f); err != nil {
		panic(err)
	}
}

func removeMappedFolderWithTui(name string) {
	c, err := winbox.Load(name + winbox.Ext)
	if err != nil {
		panic(err)
	}

	var mapped []string
	for _, mf := range c.MappedFolders {
		mapped = append(mapped, mf.HostFolder)
	}

	r := clist.Show("Remove Mapped Folder", "Removing", mapped)
	c.RemoveMappedFolder(r)

	f, err := os.Create(name + winbox.Ext)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := c.WriteXML(f); err != nil {
		panic(err)
	}
}

func isReadOnly(s string) bool {
	switch s {
	case "yes", "y", "true", "t":
		return true
	default:
		return false
	}
}

func replaceEnv(s string) string {
	e := extractEnv(s)
	v := resolveEnv(e)
	return strings.ReplaceAll(s, e, v)
}

func indexesOf(s string, ch rune) (int, int) {
	var idx []int
	for i, c := range s {
		if c == ch {
			idx = append(idx, i)
		}
	}
	if len(idx) != 2 {
		return 0, len(s) - 1
	}
	return idx[0], idx[1] + 1
}

func extractOldEnv(s string) string {
	_s, e := indexesOf(s, '%')
	return s[_s:e]
}

func replaceOldEnv(s string) string {
	e := extractOldEnv(s)
	v := os.Getenv(strings.ReplaceAll(e, "%", ""))
	if v == "" {
		return s
	}
	return strings.ReplaceAll(s, e, v)
}

func resolveEnv(s string) string {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		log.Printf("invalid env variable format: %s\nexpected: $env:Name", s)
		return ""
	}
	if parts[0] != "$env" {
		log.Printf("invalid env variable format; %s\nexpected: $env:Name", s)
	}
	return os.Getenv(parts[1])
}

func extractEnv(s string) string {
	return regexp.MustCompile(`\$env:[^\\]+`).FindString(s)
}

func init() {
	mapFolder.Flags().StringP("name", "N", "sandbox", "Name of the Windows Sandbox configuration")
	mapFolder.Flags().StringP("host", "H", "", "Host folder")
	mapFolder.Flags().StringP("sandbox", "S", "", "Sandbox folder")
	mapFolder.Flags().BoolP("readonly", "R", false, "Read-only")
	mapFolder.Flags().BoolP("tui", "u", false, "Use the TUI to map a folder")
	mapFolder.Flags().BoolP("remove", "r", false, "Remove a mapped folder")

	rootCmd.AddCommand(mapFolder)
}
