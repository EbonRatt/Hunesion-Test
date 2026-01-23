# Start Docker Services Script (Windows PowerShell)
# This script starts Kafka and Ansible containers

Write-Host "Starting Docker services..." -ForegroundColor Green

# Navigate to project root
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptPath
Set-Location $projectRoot

# Check if Docker is running
try {
    docker ps | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error: Docker is not running. Please start Docker Desktop." -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error: Docker is not installed or not running." -ForegroundColor Red
    exit 1
}

# Start services
Write-Host "Starting Ansible container..." -ForegroundColor Yellow
docker compose up -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "Services started successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Running containers:" -ForegroundColor Cyan
    docker compose ps
} else {
    Write-Host "Error: Failed to start services." -ForegroundColor Red
    exit 1
}
