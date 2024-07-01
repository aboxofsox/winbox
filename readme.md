# winbox

A small tool to manage [Windows Sandbox]('https://learn.microsoft.com/en-us/windows/security/application-security/application-isolation/windows-sandbox/windows-sandbox-overview') configurations.

## Basic Usage

```shell
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

Flags:
  -h, --help     help for winbox
  -t, --toggle   Help message for toggle

Use "winbox [command] --help" for more information about a command.

```

## `create`

```powershell
Create a new Windows Sandbox configuration file

Usage:
  winbox create [flags]

Flags:
  -a, --audio string        Audio input (default "Disable")
  -c, --clipboard string    Clipboard redirection (default "Disable")
  -h, --help                help for create
  -m, --memory int          Memory in MB (default 8192)
  -N, --name string         Name of the Windows Sandbox configuration (default "sandbox")
  -n, --networking string   Networking configuration (default "Default")
  -r, --printer string      Printer redirection (default "Disable")
  -p, --protected string    Protected client (default "Disable")
  -v, --vGpu string         Enable or disable vGPU (default "Disable")
  -i, --video string        Video input (default "Disable")
```

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
