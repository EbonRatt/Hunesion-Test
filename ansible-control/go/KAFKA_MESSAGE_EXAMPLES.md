# Kafka Message Structure Examples

This document provides examples of Kafka message structures for various operations.

## Enable User (Re-enable after soft delete)

### Message Structure

Send this message to the **consumer topic** (configured in `config.json`):

```json
{
  "event-type": "UPDATE",
  "playbook": "enable_user",
  "target-server": ["server1", "server2"],
  "OS-type": "Linux",
  "user": "admin",
  "password": "your-password",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "username": "john",
    "shell": "/bin/bash"
  }
}
```

### Required Fields

- **`event-type`**: Must be `"UPDATE"` (case-insensitive)
- **`playbook`**: Must be `"enable_user"`
- **`target-server`**: Array of server hostnames/IPs from your inventory (e.g., `["server1", "server2"]`)
- **`payload`**: JSON object containing:
  - **`username`** (required): The username to enable
  - **`shell`** (optional): Shell to restore (defaults to `/bin/bash` if not provided)

### Optional Fields

- **`OS-type`**: Operating system type (e.g., `"Linux"`)
- **`user`**: SSH user for connection
- **`password`**: SSH password
- **`timestamp`**: Timestamp string

### Minimal Example

```json
{
  "event-type": "UPDATE",
  "playbook": "enable_user",
  "target-server": ["server1"],
  "payload": {
    "username": "john"
  }
}
```

### Example with Custom Shell

```json
{
  "event-type": "UPDATE",
  "playbook": "enable_user",
  "target-server": ["server1", "server2"],
  "payload": {
    "username": "john",
    "shell": "/bin/zsh"
  }
}
```

---

## Other Operation Examples

### Create User

```json
{
  "event-type": "CREATE",
  "playbook": "create_user",
  "target-server": ["server1", "server2"],
  "OS-type": "Linux",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "username": "john",
    "user_shell": "/bin/bash",
    "user_groups": ["users", "sudo"],
    "plain_password": "secure-password-123"
  }
}
```

### Delete User (Soft Delete - Disable)

```json
{
  "event-type": "DELETE",
  "playbook": "delete_user",
  "target-server": ["server1", "server2"],
  "OS-type": "Linux",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "username": "john"
  }
}
```

### Change Password

```json
{
  "event-type": "UPDATE",
  "playbook": "change_password",
  "target-server": ["server1", "server2"],
  "OS-type": "Linux",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "username": "john",
    "plain_password": "new-password-123"
  }
}
```

### Add Server to Inventory

```json
{
  "event-type": "CREATE",
  "playbook": "add_server",
  "target-server": [],
  "OS-type": "Linux",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "alias": "server1",
    "ansible_host": "192.168.1.100",
    "ansible_user": "root",
    "ansible_port": 22,
    "ansible_password": "ssh-password",
    "ansible_become_method": "sudo",
    "ansible_become_password": "sudo-password",
    "group": "targets"
  }
}
```

### Remove Server from Inventory

```json
{
  "event-type": "DELETE",
  "playbook": "remove_server",
  "target-server": [],
  "OS-type": "Linux",
  "timestamp": "2024-01-15T10:30:00Z",
  "payload": {
    "alias": "server1"
  }
}
```

---

## Response Messages

After processing, the system sends a response to the **producer topic** (configured in `config.json`):

### Success Response

```json
{
  "status": "success",
  "event-type": "UPDATE",
  "playbook": "enable_user",
  "target-server": ["server1", "server2"],
  "message": "Successfully enabled user 'john' on target server(s)",
  "timestamp": "2024-01-15T10:30:05Z"
}
```

### Error Response

```json
{
  "status": "error",
  "event-type": "UPDATE",
  "playbook": "enable_user",
  "target-server": ["server1", "server2"],
  "message": "Failed to process UPDATE event",
  "timestamp": "2024-01-15T10:30:05Z",
  "error": "User 'john' does not exist on target server(s)"
}
```

---

## Configuration

Kafka topics are configured in `config/config.json`:

```json
{
  "kafka": {
    "broker": "localhost:29092",
    "consumer_topic": "consumer-topic",
    "producer_topic": "producer-topic"
  }
}
```

- **`consumer_topic`**: Topic where you send commands/events
- **`producer_topic`**: Topic where responses are sent

---

## Notes

1. **Server Names**: The `target-server` array must contain server aliases that exist in your Ansible inventory (`internal/inventory/files/hosts.ini`)

2. **File-Only Mode**: If `target-server` is empty or playbook name contains `_file`, the playbook will only create the vars file without execution

3. **Password Security**: After successful operations, password fields in vars files are cleared for security

4. **Error Handling**: The system detects common errors:
   - User not found
   - Server not found in inventory
   - Playbook execution failures
