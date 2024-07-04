package winbox

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	WindowsSandboxPath string `json:"windowsSandboxPath"`
}

type Configuration struct {
	VGpu                 string         `xml:"VGpu" json:"vGpu"`
	Networking           string         `xml:"Networking" json:"networking"`
	MappedFolders        []MappedFolder `xml:"MappedFolders>MappedFolder" json:"mappedFolders"`
	LogonCommand         Command        `xml:"LogonCommand>Command" json:"logonCommand"`
	AudioInput           string         `xml:"AudioInput" json:"audioInput"`
	VideoInput           string         `xml:"VideoInput" json:"videoInput"`
	ProtectedClient      string         `xml:"ProtectedClient" json:"protectedClient"`
	PrinterRedirection   string         `xml:"PrinterRedirection" json:"printerRedirection"`
	ClipboardRedirection string         `xml:"ClipboardRedirection" json:"clipboardRedirection"`
	MemoryInMB           int            `xml:"MemoryInMB" json:"memoryInMB"`
}

type MappedFolder struct {
	HostFolder    string `xml:"HostFolder" json:"hostFolder"`
	SandboxFolder string `xml:"SandboxFolder" json:"sandboxFolder"`
	ReadOnly      bool   `xml:"ReadOnly" json:"readOnly"`
}

type Command struct {
	Command string `xml:",chardata" json:"command"`
}

const Ext = ".wsb"

// WindowsSandboxPath is the default path where Windows Sandbox is installed
// If that is not where it is installed, you can create a 'config.json' file
// and define a custom path. The file should be in the same directory as the
// executable and have the following format:
//
//	{
//	  "windowsSandboxPath": "C:\\Windows\\System32\\WindowsSandbox.exe"
//	}
var WindowsSandboxPath = "C:\\Windows\\System32\\WindowsSandbox.exe"

// Load loads a configuration from a file
func Load(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch filepath.Ext(path) {
	case ".json":
		return DecodeJSON(f)
	case ".xml", ".wsb":
		return DecodeXML(f)
	default:
		return nil, fmt.Errorf("unsupported file type %s", filepath.Ext(path))
	}
}

func LoadWinboxConfig() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := new(Config)
	dec := json.NewDecoder(f)
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil

}

// WriteXML writes the configuration to an XML file
func (c *Configuration) WriteXML(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(c)
}

// XML converts the configuration to XML
func (c *Configuration) XML() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := c.WriteXML(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteJSON writes the configuration to a JSON file
func (c *Configuration) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(c)
}

// JSON converts the configuration to JSON
func (c *Configuration) JSON() ([]byte, error) {
	return json.Marshal(c)
}

// AddMappedFolder adds a mapped folder on the host to Windows Sandbox
func (c *Configuration) AddMappedFolder(mf MappedFolder) {
	c.MappedFolders = append(c.MappedFolders, mf)
}

// AddLogonCommand adds a command to run on startup
func (c *Configuration) AddLogonCommand(cmd Command) {
	c.LogonCommand = cmd
}

// DecodeJSON decodes a JSON configuration
func DecodeJSON(r io.Reader) (*Configuration, error) {
	c := new(Configuration)
	dec := json.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

// DecodeXML decodes an XML configuration
func DecodeXML(r io.Reader) (*Configuration, error) {
	c := new(Configuration)
	dec := xml.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

func Run(name string) error {
	cmd := exec.Command("cmd", "/c", "start", name+Ext)
	return cmd.Run()
}
