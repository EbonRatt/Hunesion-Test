# Run Agent Script (Windows PowerShell)
# This script runs the ansible-control agent

Write-Host "Starting Ansible Control Agent..." -ForegroundColor Green

# Navigate to project root
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptPath
$goDir = Join-Path $projectRoot "ansible-control\go"

# Check for executable in bin directory first (recommended), then in go directory
$binDir = Join-Path $projectRoot "ansible-control\bin"
$exePath = Join-Path $binDir "ansible-control.exe"
if (-not (Test-Path $exePath)) {
    # Fallback to go directory
    $exePath = Join-Path $goDir "ansible-control.exe"
    if (-not (Test-Path $exePath)) {
        Write-Host "Error: ansible-control.exe not found!" -ForegroundColor Red
        Write-Host "Please build the application first:" -ForegroundColor Yellow
        Write-Host "  cd ansible-control\go" -ForegroundColor Yellow
        Write-Host "  mkdir ..\bin" -ForegroundColor Yellow
        Write-Host "  go build -o ..\bin\ansible-control.exe main.go" -ForegroundColor Yellow
        exit 1
    }
}

# Check if Docker services are running
Write-Host "Checking Docker services..." -ForegroundColor Yellow
$ansibleContainer = docker ps --filter "name=ansible" --format "{{.Names}}"
if (-not $ansibleContainer) {
    Write-Host "Warning: Ansible container is not running." -ForegroundColor Yellow
    Write-Host "Starting services..." -ForegroundColor Yellow
    & "$scriptPath\start-services.ps1"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Failed to start services." -ForegroundColor Red
        exit 1
    }
}

# Change to executable directory
$exeDir = Split-Path -Parent $exePath
Set-Location $exeDir

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Ansible Control Agent" -ForegroundColor Cyan
Write-Host "  Main File: main.go" -ForegroundColor Cyan
Write-Host "  Executable: ansible-control.exe" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Run the agent
& $exePath
