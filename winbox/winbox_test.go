package winbox

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWinbox(t *testing.T) {
	t.Run("TestXML", testXml)
}

func testXml(t *testing.T) {
	c := &Configuration{}
	c.VGpu = "Disable"
	c.Networking = "Default"
	c.AudioInput = "Disable"
	c.VideoInput = "Disable"
	c.ProtectedClient = "Disable"
	c.PrinterRedirection = "Disable"
	c.ClipboardRedirection = "Disable"
	c.MemoryInMB = 8 * 1024

	mf := MappedFolder{
		HostFolder:    filepath.Join("C:", "Users", "user", "Desktop"),
		SandboxFolder: filepath.Join("C:", "Users", "WDAGUtilityAccount", "Desktop"),
	}

	c.MappedFolders = append(c.MappedFolders, mf)

	lc := Command{
		Command: "notepad.exe",
	}

	c.LogonCommand = lc

	f, err := os.Create("test.wsb")
	if err != nil {
		t.Fatal(err)
	}

	if err := c.WriteXML(f); err != nil {
		t.Fatal(err)
	}
	f.Close()

	nc, err := Load("test.wsb")
	if err != nil {
		t.Fatal(err)
	}

	if c.VGpu != nc.VGpu {
		t.Fatalf("expected %s, got %s", c.VGpu, nc.VGpu)
	}
	if c.Networking != nc.Networking {
		t.Fatalf("expected %s, got %s", c.Networking, nc.Networking)
	}
	if c.AudioInput != nc.AudioInput {
		t.Fatalf("expected %s, got %s", c.AudioInput, nc.AudioInput)
	}
	if c.VideoInput != nc.VideoInput {
		t.Fatalf("expected %s, got %s", c.VideoInput, nc.VideoInput)
	}
	if c.ProtectedClient != nc.ProtectedClient {
		t.Fatalf("expected %s, got %s", c.ProtectedClient, nc.ProtectedClient)
	}
	if c.PrinterRedirection != nc.PrinterRedirection {
		t.Fatalf("expected %s, got %s", c.PrinterRedirection, nc.PrinterRedirection)
	}
	if c.ClipboardRedirection != nc.ClipboardRedirection {
		t.Fatalf("expected %s, got %s", c.ClipboardRedirection, nc.ClipboardRedirection)
	}
	if c.MemoryInMB != nc.MemoryInMB {
		t.Fatalf("expected %d, got %d", c.MemoryInMB, nc.MemoryInMB)
	}
	if c.MappedFolders[0].HostFolder != nc.MappedFolders[0].HostFolder {
		t.Fatalf("expected %s, got %s", c.MappedFolders[0].HostFolder, nc.MappedFolders[0].HostFolder)
	}
	if c.MappedFolders[0].SandboxFolder != nc.MappedFolders[0].SandboxFolder {
		t.Fatalf("expected %s, got %s", c.MappedFolders[0].SandboxFolder, nc.MappedFolders[0].SandboxFolder)
	}
	if c.MappedFolders[0].ReadOnly != nc.MappedFolders[0].ReadOnly {
		t.Fatalf("expected %t, got %t", c.MappedFolders[0].ReadOnly, nc.MappedFolders[0].ReadOnly)
	}
	if c.LogonCommand.Command != nc.LogonCommand.Command {
		t.Fatalf("expected %s, got %s", c.LogonCommand.Command, nc.LogonCommand.Command)
	}

	os.Remove("test.wsb")

}
