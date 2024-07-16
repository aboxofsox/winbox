param(
        [string] $InstallPath = (Join-Path -Path $env:USERPROFILE -ChildPath 'winbox'), 
        [version] $Version = '0.3.0',
        [boolean] $AddToPath = $false
)

function Install-Winbox {
        [CmdletBinding()]
        param(
                [Parameter(Mandatory=$false)]
                [string] $InstallPath = (Join-Path -Path $env:USERPROFILE -ChildPath 'winbox'),
                [Parameter(Mandatory=$false)]
                [boolean] $AddToPath = $false
        )

        $arch = ($env:PROCESSOR_ARCHITECTURE).ToLower()
        Write-Log -Message "Detected architecture: $arch"

        $baseURL = "https://github.com/aboxofsox/winbox/releases/download/$($version.ToString())/winbox-$arch-v$($version.ToString()).exe"
        $tempDir = Join-Path -Path $env:TEMP -ChildPath (New-Guid).Guid
        $tempPath = Join-Path -Path $tempDir -ChildPath 'winbox.exe'

        if (-not (Test-Path -Path $InstallPath)) {
                New-Item -ItemType directory -Path $InstallPath -Force | Out-Null
        }

        if (-not (Test-Path -Path $tempDir)) {
                New-Item -ItemType directory -Path $tempDir -Force | Out-Null
        }

        Write-Log -Message "Download winbox from $baseURL to $tempPath"
        Invoke-Download -URL $baseURL -OutputPath $tempPath

        $checksum = Get-FileHash -Path $tempPath -Algorithm SHA256
        Write-Log -Message "Checksum: $($checksum.Hash)"

        # download checksum file
        # https://github.com/aboxofsox/winbox/releases/download/0.3.0/checksums.txt
        $checksumURL = "https://github.com/aboxofsox/winbox/releases/download/$($version.ToString())/checksums.txt"
        $checksumPath = Join-Path -Path $tempDir -ChildPath 'checksums.txt'
        Invoke-Download -URL $checksumURL -OutputPath $checksumPath

        $content = Get-Content -Path $checksumPath
        $match = $false
        foreach ($line in $content) {
                $hash = $line -split ' '
                if ($hash[0] -eq $checksum.Hash) {
                        $match = $true
                        Write-Log -Message "Checksum matched: $($hash[0])"
                } 
        }
        if (-not $match) {
                Write-Log -Message "Checksum mismatch: $($checksum.Hash)"
                Write-Host "❌ Checksum mismatch: $($checksum.Hash)" -ForegroundColor Red
                exit
        }

        Write-Log -Message "Copy winbox to $InstallPath"
        Copy-Item -Path (Join-Path -Path $tempDir -ChildPath 'winbox.exe') -Destination $InstallPath -Recurse -Force

        if ($AddToPath) {
                $path = (Get-Env -Key 'Path') -split ';'
                if ($path -notcontains $InstallPath) {
                        $path += $InstallPath
                        Write-Env -Key 'Path' -Value ($path -join ';')
                        $env:PATH = $path
                } else {
                        Write-Log -Message 'Winbox is already in PATH. Skipping'
                }
        }

        Remove-Item -Path $tempDir -Recurse -Force

        Write-Log -Message 'Set environment variables'
        Write-Log -Message 'WINBOX_USER=WDAGUtilityAccount'
        Write-Env -Key 'WINBOX_USER' -Value 'WDAGUtilityAccount'
        Write-Log -Message "WINBOX_DIR=$InstallPath"
        Write-Env -Key 'WINBOX_DIR' -Value $InstallPath

        Write-Log -Message "Winbox has finished installation"
        Write-Host "✅ Winbox has finished installation" -ForegroundColor Green
        Write-Host "📃 You can view the log in $(Join-Path -Path $InstallPath -ChildPath 'logs')"
        Write-Host "💡 Added environment variables:`n`t> WINBOX_USER`n`t> WINBOX_DIR" -ForegroundColor Yellow
}

function Invoke-Download {
        [CmdletBinding()]
        param(
                [Parameter(Mandatory=$true)]
                [string] $URL,
                [Parameter(Mandatory=$true)]
                [string] $OutputPath
        )

        $webClient = New-Object System.Net.WebClient
        $webClient.DownloadFile($URL, $OutputPath)

        if (-not (Test-Path -Path $OutputPath)) {
                Write-Log -Message "Failed to download $URL to $OutputPath"
                throw "Failed to download $URL to $OutputPath"
        }
}

function Write-Path {
        [CmdletBinding()]
        param(
                [Parameter(Mandatory=$true)]
                [string] $Path
        )

        $path = [System.Environment]::GetEnvironmentVariable('PATH', 'User')
        $path = $path -split ';'
        if ($path -notcontains $Path) {
                $path += $Path
                $path = $path -join ';'
                [System.Environment]::SetEnvironmentVariable('PATH', $path, 'User')
        }
}

function Write-Log {
        [CmdletBinding()]
        param(
                [Parameter()]
                [string] $Message
        )

        $logdir = Join-Path -Path $InstallPath -ChildPath 'logs'
        $logpath = Join-Path -Path $logdir -ChildPath 'winbox-install.log'

        if (-not (Test-Path -Path $logdir)) {
                New-Item -ItemType directory -Path $logdir -Force | Out-Null
        }

        Write-Host "> $Message"
        Format-LogMessage -Message $Message | Out-File -FilePath $logpath -Append
}

function Format-LogMessage {
        [CmdletBinding()]
        param(
                [Parameter(Mandatory=$true)]
                [string] $Message
        )

        return "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') - $Message"

}

function Publish-Env {
        if (-not ("Win32.NativeMethods" -as [Type])) {
                Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
        uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
        }
        $HWND_BROADCAST = [IntPtr] 0xffff
        $WM_SETTINGCHANGE = 0x1a
        $result = [UIntPtr]::Zero
        [Win32.NativeMethods]::SendMessageTimeout($HWND_BROADCAST,
                $WM_SETTINGCHANGE,
                [UIntPtr]::Zero,
                "Environment",
                2,
                5000,
                [ref] $result
        ) | Out-Null
}

function Write-Env {
        param(
                [Parameter(Mandatory = $true)]
                [string]$Key,
                
                [string]$Value
        )

        $RegisterKey = Get-Item -Path 'HKCU:'
        $EnvRegisterKey = $RegisterKey.OpenSubKey('Environment', $true)
        
        if ($null -eq $Value) {
                $EnvRegisterKey.DeleteValue($Key)
        } else {
                $RegistryValueKind = if ($Value.Contains('%')) {
                        [Microsoft.Win32.RegistryValueKind]::ExpandString
                } elseif ($EnvRegisterKey.GetValue($Key)) {
                        $EnvRegisterKey.GetValueKind($Key)
                } else {
                        [Microsoft.Win32.RegistryValueKind]::String
                }
                $EnvRegisterKey.SetValue($Key, $Value, $RegistryValueKind)
        }

        Publish-Env
}

function Get-Env {
        param([string]$Key)
        $regkey = Get-Item 'HKCU:'
        $envkey = $regkey.OpenSubKey('Environment')
        $envkey.GetValue($Key, $null, [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
}

Install-Winbox -InstallPath $InstallPath -AddToPath $AddToPath