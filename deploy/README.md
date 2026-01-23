# Deployment Guide

## Main Entry Point

**The main file to run the agent is:**
```
ansible-control/go/main.go
```

**After building, the executable should be in:**
- Windows: `ansible-control/bin/ansible-control.exe` (recommended)
- Linux: `ansible-control/bin/ansible-control` (recommended)
- macOS: `ansible-control/bin/ansible-control` (recommended)

**Note:** Build to `bin/` directory to keep source code clean. The scripts will check `bin/` first, then fallback to `go/` directory.

## Quick Start

1. **Start Docker services (Kafka, Ansible container):**
   ```bash
   .\start-services.ps1    # Windows PowerShell
   # or
   ./start-services.sh     # Linux/macOS
   ```

2. **Run the agent:**
   ```bash
   .\run-agent.ps1         # Windows PowerShell
   # or
   ./run-agent.sh          # Linux/macOS
   ```

3. **Stop services:**
   ```bash
   .\stop-services.ps1     # Windows PowerShell
   # or
   ./stop-services.sh      # Linux/macOS
   ```

## File Structure

```
Ansible/
├── deploy/                    # Deployment scripts
│   ├── README.md             # This file
│   ├── start-services.ps1    # Start Docker services (Windows)
│   ├── start-services.sh     # Start Docker services (Linux/macOS)
│   ├── stop-services.ps1     # Stop Docker services (Windows)
│   ├── stop-services.sh      # Stop Docker services (Linux/macOS)
│   └── run-agent.ps1         # Run the agent (Windows)
│   └── run-agent.sh          # Run the agent (Linux/macOS)
├── docker-compose.yml         # Docker Compose configuration
├── ansible-control/
│   ├── bin/                  # ⭐ Build output directory (executables go here)
│   │   └── ansible-control.exe
│   └── go/
│       └── main.go           # ⭐ MAIN ENTRY POINT (source code)
└── README.md
```

## Prerequisites

- Docker Desktop (Windows) or Docker Engine (Linux/macOS)
- Docker Compose
- Go 1.21+ (for building, not required for running the executable)
- **Built executable** (`ansible-control.exe` or `ansible-control`)
- **Configuration file** (`ansible-control/go/config/config.json`)

## Configuration

**Before running, configure Kafka settings in:**
```
ansible-control/go/config/config.json
```

**Edit the file to change:**
- Kafka broker address
- Consumer topic (where events are received)
- Producer topic (where responses are sent)

See `ansible-control/go/config/README.md` for details.

## ⚠️ Important Notes

**The deploy folder does NOT build the executable for you!**

**Before using deploy scripts, you must:**
1. Build the executable first:
   ```powershell
   cd ansible-control\go
   go build -o ansible-control.exe main.go
   ```

2. Then use deploy scripts to run it:
   ```powershell
   cd deploy
   .\start-services.ps1
   .\run-agent.ps1
   ```

**The deploy folder automates running, not building!**

## Important Notes

⚠️ **The deploy folder does NOT build the executable for you!**

**Before using deploy scripts, you must:**
1. Build the executable first (to bin directory - recommended):
   ```powershell
   cd ansible-control\go
   mkdir ..\bin
   go build -o ..\bin\ansible-control.exe main.go
   
   # Or use the build script:
   .\build.ps1
   ```

2. Then use deploy scripts to run it:
   ```powershell
   cd deploy
   .\start-services.ps1
   .\run-agent.ps1
   ```

**The deploy folder automates running, not building!**

**Note:** The scripts will automatically look for the executable in `bin/` directory first, then fallback to `go/` directory.