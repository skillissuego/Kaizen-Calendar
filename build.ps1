# Build script for Windows (PowerShell).
# Requires a MinGW-w64 GCC toolchain (e.g. WinLibs) with its bin directory
# on PATH, so that Fyne's CGO/OpenGL driver can compile.

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

$env:CGO_ENABLED = "1"

# Regenerate Windows resources (icon + version info) if go-winres is installed.
$winres = Join-Path (go env GOPATH) "bin\go-winres.exe"
Push-Location (Join-Path $PSScriptRoot "cmd\kaizen")
try {
    if (Test-Path $winres) {
        & $winres make
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    } else {
        Write-Warning "go-winres not found; skipping resource generation."
    }
}
finally {
    Pop-Location
}

go test ./...
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go build -ldflags "-H=windowsgui" -o kaizen.exe ./cmd/kaizen
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Built .\kaizen.exe"
