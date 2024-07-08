$packageName = 'winbox'
[version]$version = '0.3.0'
$commit = (git rev-parse --short HEAD)
$date = (Get-Date -Format 'yyyy-MM-dd')
$archs = 'windows/386', 'windows/arm64', 'windows/amd64'
$previousBuilds = Get-ChildItem '.\bin' -Filter "$packageName-*" -Recurse
$oldVersionPath = Join-Path -Path '.\old' -ChildPath $version.ToString()

if (-not (Test-Path -Path $oldVersionPath)) {
        New-Item -ItemType directory -Path $oldVersionPath -Force | Out-Null
}
$backups = Get-ChildItem $oldVersionPath -Filter '*.bak' -Recurse

foreach ($build in $previousBuilds) {
        $fileName = Split-Path $build.FullName -Leaf
        $destination = Join-Path -Path $oldVersionPath -ChildPath $fileName
        Move-Item -Path $build.FullName -Destination $destination -Force
    }

foreach ($backup in $backups) {
        if ($backup.LastWriteTime -lt (Get-Date).AddDays(30)) {
                Remove-Item $backup.FullName -Force
        }
}

go mod tidy

foreach ($arch in $archs) {
        $arch = $arch.Split('/')[1]
        $fullpath = ".\bin\$packageName-$arch-$($version.ToString()).exe"
        if (Test-Path $fullpath) {
                Write-Host "Build for $arch version $($version.ToString()) already exists.`nPlease update the version." -ForegroundColor Yellow
                exit
        }
        $env:GOOS = 'windows'
        $env:GOARCH = $arch
        go build -o "bin/$packageName-$arch-v$($version.ToString()).exe" -ldflags "-X main.version=$version -X main.commit=$commit -X main.date=$date" -v
}

$checksumFile = '.\bin\checksums.txt'
if (-not (Test-Path -Path $checksumFile)) {
        New-Item -ItemType file -Path $checksumFile -Force | Out-Null
}

$checksums = @()
Get-ChildItem -Path '.\bin' -Filter "$packageName-*.exe" | ForEach-Object {
        $checksum = Get-FileHash -Path $_.FullName -Algorithm SHA256
        $checksums += $checksum.Hash + ' ' + $_.Name
}

$checksums | Set-Content -Path $checksumFile
