# Complete Deployment Guide

## What the Deploy Folder Does

The `deploy/` folder contains **automation scripts** that make it easier to:
- ✅ Start Docker services (Ansible container)
- ✅ Run the agent
- ✅ Stop services

## What You Still Need

### 1. **Build the Executable (One-time or after code changes)**

**Before using deploy scripts, you must build the executable:**

```powershell
cd ansible-control\go
go build -o ansible-control.exe main.go
```

**Or use the build script:**
```powershell
cd ansible-control\go
.\build.ps1    # (if you create one)
```

### 2. **Prerequisites (Must be installed)**

- ✅ Docker Desktop (Windows) or Docker Engine (Linux/macOS)
- ✅ Docker Compose
- ✅ Kafka (can be in separate docker-compose or external service)
- ✅ Go 1.21+ (only needed for building, not for running)

### 3. **Project Files (Must exist)**

- ✅ `ansible-control/go/main.go` - Source code
- ✅ `ansible-control/go/ansible-control.exe` - Built executable
- ✅ `docker-compose.yml` - Docker configuration
- ✅ `ansible-control/go/internal/` - All internal packages

## Complete Workflow

### **First Time Setup:**

1. **Build the executable:**
   ```powershell
   cd ansible-control\go
   go build -o ansible-control.exe main.go
   ```

2. **Use deploy scripts:**
   ```powershell
   cd deploy
   .\start-services.ps1    # Start Docker
   .\run-agent.ps1         # Run agent
   ```

### **Daily Usage (After First Setup):**

**Just use deploy scripts:**
```powershell
cd deploy
.\start-services.ps1    # Start services
.\run-agent.ps1         # Run agent
.\stop-services.ps1     # Stop when done
```

## What Deploy Folder CANNOT Do

❌ **Cannot build the executable** - You must build it manually first  
❌ **Cannot install Docker** - Must be installed separately  
❌ **Cannot install Kafka** - Must be running separately  
❌ **Cannot create source code** - Code must exist in `ansible-control/go/`

## What Deploy Folder CAN Do

✅ **Start/Stop Docker services** - Automates `docker compose up/down`  
✅ **Run the agent** - Automates running `ansible-control.exe`  
✅ **Check prerequisites** - Verifies Docker and executable exist  
✅ **Auto-start services** - Starts Docker if not running when agent starts

## Summary

**Deploy folder = Automation scripts for running, NOT a complete standalone solution**

**You need:**
1. Built executable (`ansible-control.exe`)
2. Docker installed and running
3. Kafka running
4. Source code in `ansible-control/go/`

**Then deploy scripts make it easy to:**
- Start everything
- Run the agent
- Stop everything
