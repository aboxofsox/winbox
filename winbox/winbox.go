package winbox

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
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
	HostFolder    string `xml:"HostFolder"`
	SandboxFolder string `xml:"SandboxFolder"`
	ReadOnly      bool   `xml:"ReadOnly"`
}

type Command struct {
	Command string `xml:",chardata"`
}

const Ext = ".wsb"

var WindowsSandboxPath = "C:\\Windows\\System32\\WindowsSandbox.exe"

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

func (c *Configuration) Run(name string) error {

	return nil
}

func (c *Configuration) WriteXML(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(c)
}

func (c *Configuration) XML() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := c.WriteXML(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Configuration) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(c)
}

func (c *Configuration) JSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Configuration) AddMappedFolder(mf MappedFolder) {
	c.MappedFolders = append(c.MappedFolders, mf)
}

func (c *Configuration) AddLogonCommand(cmd Command) {
	c.LogonCommand = cmd
}

func DecodeJSON(r io.Reader) (*Configuration, error) {
	c := new(Configuration)
	dec := json.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}

func DecodeXML(r io.Reader) (*Configuration, error) {
	c := new(Configuration)
	dec := xml.NewDecoder(r)
	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
