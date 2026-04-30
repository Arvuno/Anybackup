param(
    [string]$SourceDir = "E:\Code\Service\foundation_cli",
    [string]$OutputPath = ""
)

$ErrorActionPreference = "Stop"

$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSCommandPath)
if ([string]::IsNullOrWhiteSpace($OutputPath)) {
    $OutputPath = Join-Path $RepoRoot "bin\foundation-cli"
}

if (-not (Test-Path -LiteralPath $SourceDir -PathType Container)) {
    throw "foundation_cli source directory not found: $SourceDir"
}

$Go = Get-Command go -ErrorAction SilentlyContinue
if (-not $Go) {
    throw "go was not found in PATH. Install Go 1.22+ or add go.exe to PATH before building foundation-cli."
}

$OutputDir = Split-Path -Parent $OutputPath
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
if (Test-Path -LiteralPath $OutputPath) {
    Remove-Item -LiteralPath $OutputPath -Force
}

$previousGoos = $env:GOOS
$previousGoarch = $env:GOARCH
$previousCgo = $env:CGO_ENABLED

try {
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"

    Push-Location $SourceDir
    try {
        & $Go.Source build -trimpath -ldflags "-s -w" -o $OutputPath ".\cmd\foundation-cli"
        if ($LASTEXITCODE -ne 0) {
            throw "go build failed with exit code $LASTEXITCODE"
        }
    }
    finally {
        Pop-Location
    }
}
finally {
    $env:GOOS = $previousGoos
    $env:GOARCH = $previousGoarch
    $env:CGO_ENABLED = $previousCgo
}

$artifact = Get-Item -LiteralPath $OutputPath
if ($artifact.Length -lt 1048576) {
    throw "built foundation-cli is unexpectedly small ($($artifact.Length) bytes): $OutputPath"
}

$hash = Get-FileHash -Algorithm SHA256 -LiteralPath $OutputPath
Write-Host "Built Linux foundation-cli: $OutputPath"
Write-Host "Size: $($artifact.Length) bytes"
Write-Host "SHA256: $($hash.Hash)"
