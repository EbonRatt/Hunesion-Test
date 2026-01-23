# Stop Docker Services Script (Windows PowerShell)

Write-Host "Stopping Docker services..." -ForegroundColor Yellow

# Navigate to project root
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptPath
Set-Location $projectRoot

# Stop services
docker compose down

if ($LASTEXITCODE -eq 0) {
    Write-Host "Services stopped successfully!" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to stop services." -ForegroundColor Red
    exit 1
}
