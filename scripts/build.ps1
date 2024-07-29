param(
        [string] $PackageName = 'winbox',
        [version] $Version = (Get-Content -Path '.\VERSION' -Raw),
        [string] $Commit = (git rev-parse --short HEAD),
        [string[]] $Archs = @('windows/386', 'windows/arm64', 'windows/amd64'),
        [switch]$TestBuild = $false
)

if (-not (Test-Path -Path 'VERSION'))  {
        Write-Host 'VERSION file not found' -ForegroundColor Red
        exit
}

$date = (Get-Date -Format 'yyyy-MM-dd')
$bin = '.\bin'
if ($TestBuild) {
        $bin = '.\testing'
}

if (-not (Test-Path -Path $bin)) {
        New-Item -ItemType directory -Path $bin -Force | Out-Null
}

function Invoke-Command {
        param(
                [string]$Command,
                [string]$Arguments
        )
        $process = Start-Process -FilePath $Command -ArgumentList $Arguments -PassThru -Wait
        $process.ExitCode
}

function Invoke-Download {
        param(
                [string]$Url,
                [string]$Destination
        )
        $webClient = New-Object System.Net.WebClient
        $webClient.DownloadFile($Url, $Destination)
}

go mod tidy

foreach ($arch in $Archs) {
        $arch = $arch.Split('/')[1]
        $fullpath = ".\bin\$PackageName-$arch-$($Version.ToString()).exe"
        if (Test-Path $fullpath) {
                Write-Host "Build for $arch version $($Version.ToString()) already exists.`nPlease update the version." -ForegroundColor Yellow
                exit
        }
        $env:GOOS = 'windows'
        $env:GOARCH = $arch
        go build -o ".\bin\$PackageName-$arch-v$($Version.ToString()).exe" -ldflags "-X main.version=$Version -X main.commit=$Commit -X main.date=$date" -v
}

$checksumFile = '.\bin\checksums.txt'
if (-not (Test-Path -Path $checksumFile)) {
        New-Item -ItemType file -Path $checksumFile -Force | Out-Null
}

$checksums = @()
Get-ChildItem -Path $bin -Filter "$PackageName-*.exe" | ForEach-Object {
        $checksum = Get-FileHash -Path $_.FullName -Algorithm SHA256
        $checksums += $checksum.Hash + ' ' + $_.Name
}

$checksums | Set-Content -Path $checksumFile
