# Build Instructions

## Important Note
**Platform-Specific Binaries:**
- `ansible-control.exe` - **Windows only** (works on Windows 10/11, Windows Server)
- `ansible-control` (Linux) - **Linux only** (works on Ubuntu, CentOS, etc.)
- `ansible-control` (macOS) - **macOS only** (works on macOS)

**You cannot use a Windows .exe file on Linux or macOS.** You need to build a separate binary for each platform.

## Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose (for running Ansible playbooks)

## Build Steps

### 1. Install Dependencies
```bash
cd ansible-control/go
go mod download
```

### 2. Build the Application

**Recommended: Build to a separate `bin/` directory to keep source code clean**

#### For Windows (current system):
```bash
# Create bin directory (one-time)
mkdir ..\bin

# Build to bin directory (recommended)
go build -o ..\bin\ansible-control.exe main.go

# Or build in current directory (not recommended - clutters source)
go build -o ansible-control.exe main.go
```

#### For Linux (from Windows PowerShell):
```powershell
# Create bin directory (one-time)
mkdir ..\bin

# Build to bin directory
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ..\bin\ansible-control-linux main.go
```

#### For macOS (from Windows PowerShell):
```powershell
# Create bin directory (one-time)
mkdir ..\bin

# Build to bin directory
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o ..\bin\ansible-control-macos main.go
```

#### For Linux (from Linux/macOS terminal):
```bash
GOOS=linux GOARCH=amd64 go build -o ansible-control main.go
```

#### For macOS (from Linux/macOS terminal):
```bash
GOOS=darwin GOARCH=amd64 go build -o ansible-control main.go
```

### 3. Run the Application

#### Windows:
```bash
.\ansible-control.exe
```

#### Linux/macOS:
```bash
./ansible-control
```

## Build Options

### Build with optimizations:
```bash
go build -ldflags="-s -w" -o ansible-control main.go
```

### Build for production (stripped binary):
```bash
go build -trimpath -ldflags="-s -w" -o ansible-control main.go
```

## Build for Multiple Platforms

### Build all platforms at once (from Windows PowerShell):
```powershell
# Windows
go build -o ansible-control-windows.exe main.go

# Linux
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o ansible-control-linux main.go

# macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o ansible-control-macos main.go
```

### Result:
- `ansible-control-windows.exe` - Use on Windows
- `ansible-control-linux` - Use on Linux
- `ansible-control-macos` - Use on macOS

## Docker Build (Optional)

If you want to containerize the Go application:

```bash
# Build Docker image
docker build -t ansible-control:latest -f Dockerfile.go .

# Run container
docker run -d ansible-control:latest
```

## Development

### Run without building:
```bash
go run main.go
```

### Run with live reload (using air or similar):
```bash
# Install air: go install github.com/cosmtrek/air@latest
air
```
