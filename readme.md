# winbox

A small tool to manage [Windows Sandbox]('https://learn.microsoft.com/en-us/windows/security/application-security/application-isolation/windows-sandbox/windows-sandbox-overview') configurations.

## Prerequisites

This tool assumes you have installed Windows Sandbox in it's default directory. You can create a `config.json` file in the same directory as the `winbox` executable if you have it installed elsewhere, somehow.

**Windows Sandbox is not available for Windows Home**

To enable Windows Sandbox:

```powershell
Enable-WindowsOptionalFeature -FeatureName 'Containers-DisposableClientVM' -All -Online
```

### Example `config.json`

```json
{
  "windowsSandboxPath": "D:\\Path\\To\\WindowsSandbox\\WindowsSandbox.exe"
}
```

## Installing

```
powershell.exe -c 'https://raw.githubusercontent.com/aboxofsox/winbox/main/scripts/install.ps1' | iex
```

If you have Go, you can simply do `go install github.com/aboxofsox/winbox`.

Otherwise, you will need to download a release, and place it somewhere on your PC. Ideally it would go somewhere that's mapped to `PATH`.

## Basic Usage

```
A small tool to manage Windows Sandbox configurations.

Usage:
  winbox [command]

Available Commands:
  add-logon   Add a logon command to the Windows Sandbox configuration
  completion  Generate the autocompletion script for the specified shell
  create      Create a new Windows Sandbox configuration
  help        Help about any command
  map         Map a folder from the host to Windows Sandbox
  run         Run a Windows Sandbox configuration
  select      Select a Windows Sandbox configuration

Flags:
  -h, --help     help for winbox
  -t, --toggle   Help message for toggle

Use "winbox [command] --help" for more information about a command.


```

## `create`

```powershell
Create a new Windows Sandbox Configuration

Usage:
  winbox create [flags]

Flags:
  -a, --audio string        Audio input (default "Disable")
  -c, --clipboard string    Clipboard redirection (default "Disable")
  -h, --help                help for create
  -m, --memory string       Memory in MB (default "1024")
  -N, --name string         Name of the Windows Sandbox configuration (default "sandbox")
  -n, --networking string   Networking configuration (default "Default")
  -r, --printer string      Printer redirection (default "Disable")
  -p, --protected string    Protected client (default "Disable")
  -u, --tui                 Use the TUI to create a configuration
  -g, --vGpu string         Enable or disable vGPU (default "Disable")
  -v, --video string        Video input (default "Disable")

```

### Expressions

When allocating memory, you can pass an expression as `Memory in Megabytes` when using the optional TUI.

_When using the optional TUI interface (`-u`), empty `name` fields will not create the default `sandbox.wsb` file._

### Example

```powershell
winbox create -N default -a Enable -m (8 * 1024)
```

## `add-logon`

```powershell
Add a logon command to the Windows Sandbox configuration

Usage:
  winbox add-logon [flags]

Flags:
  -c, --command string   Command to run on logon
  -h, --help             help for add-logon
  -N, --name string      Name of the Windows Sandbox configuration (default "sandbox")
```

### Example

```powershell
winbox add-logon -N default -c 'notepad.exe'
```

## `map`

You can pass environment variables as you would in cmd or PowerShell. (i.e. `%USERPROFILE%`, `$env:UserProfile`)

`$SandboxUser` is a special variable that holds the path to the `WDAGUtilityAccount` user folder (`C:\Users\WDAGUtilityAccount`).


```powershell
Map a folder from the host to Windows Sandbox

Usage:
  winbox map [flags]

Flags:
  -h, --help             help for map
  -H, --host string      Host folder
  -N, --name string      Name of the Windows Sandbox configuration (default "sandbox")
  -R, --readonly         Read-only
  -S, --sandbox string   Sandbox folder
```

## Example

```powershell
winbox map -N default -H "$env:UserProfile\Downloads" -S 'C:\\Users\WDAGUtilityAccount\Downloads' -R $false
```

## `run`

```powershell
Run a Windows Sandbox configuration

Usage:
  winbox run [flags]

Flags:
  -h, --help          help for run
  -N, --name string   Name of the Windows Sandbox configuration

```

### Example

```powershell
winbox run -N default
```

## `select`

`select` will open a TUI menu where you can select which Windows Sandbox configuration to use and launches Windows Sandbox.
