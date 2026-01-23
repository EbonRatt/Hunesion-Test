# Quick Start Guide

## Main Entry Point

**⭐ Main File:** `ansible-control/go/main.go`  
**⭐ Executable:** `ansible-control/go/ansible-control.exe` (Windows) or `ansible-control/go/ansible-control` (Linux/macOS)

## Deployment Steps

### 1. Start Docker Services

**Windows:**
```powershell
cd deploy
.\start-services.ps1
```

**Linux/macOS:**
```bash
cd deploy
chmod +x *.sh
./start-services.sh
```

This will:
- Start the Ansible Docker container
- Make sure all required services are running

### 2. Run the Agent

**Windows:**
```powershell
cd deploy
.\run-agent.ps1
```

**Linux/macOS:**
```bash
cd deploy
./run-agent.sh
```

This will:
- Check if the executable exists
- Verify Docker services are running
- Start the Ansible Control Agent
- The agent will listen to Kafka topic `consumer-topic` and send responses to `producer-topic`

### 3. Stop Services (when done)

**Windows:**
```powershell
cd deploy
.\stop-services.ps1
```

**Linux/macOS:**
```bash
cd deploy
./stop-services.sh
```

## Manual Run (Alternative)

If you prefer to run manually:

```bash
# 1. Start Docker services
cd D:\Hunesion\Ansible
docker compose up -d

# 2. Run the agent
cd ansible-control\go
.\ansible-control.exe    # Windows
# or
./ansible-control        # Linux/macOS
```

## File Structure

```
Ansible/
├── deploy/                          # ⭐ Deployment scripts
│   ├── README.md                    # Full documentation
│   ├── QUICK_START.md               # This file
│   ├── start-services.ps1/.sh       # Start Docker services
│   ├── stop-services.ps1/.sh        # Stop Docker services
│   └── run-agent.ps1/.sh            # Run the agent
├── docker-compose.yml                # Docker Compose config
└── ansible-control/
    └── go/
        ├── main.go                  # ⭐ MAIN ENTRY POINT
        └── ansible-control.exe       # ⭐ Built executable
```

## Troubleshooting

**Agent not starting?**
- Check if Docker is running: `docker ps`
- Check if Ansible container is running: `docker compose ps`
- Verify executable exists: `Test-Path ansible-control\go\ansible-control.exe`

**Services not starting?**
- Make sure Docker Desktop is running
- Check Docker Compose version: `docker compose version`
- View logs: `docker compose logs ansible`
