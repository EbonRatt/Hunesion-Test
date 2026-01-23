# Build Script for Windows
# Builds the executable to ../bin directory to keep source clean

Write-Host "Building Ansible Control Agent..." -ForegroundColor Green

# Create bin directory if it doesn't exist
$binDir = Join-Path (Split-Path -Parent $PSScriptRoot) "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
    Write-Host "Created bin directory: $binDir" -ForegroundColor Yellow
}

# Build for Windows
$exePath = Join-Path $binDir "ansible-control.exe"
Write-Host "Building to: $exePath" -ForegroundColor Cyan
go build -o $exePath main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Executable: $exePath" -ForegroundColor Cyan
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
