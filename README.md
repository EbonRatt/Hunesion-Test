# Ansible-like System in Go with Kafka

This is an Ansible-like orchestration system implemented in Go that listens to Kafka messages and executes playbooks on remote hosts.

## Features

- **Kafka Listener**: Consumes tasks from Kafka topics
- **Playbook System**: Predefined playbooks for common tasks
- **SSH Execution**: Remote command execution via SSH
- **User Creation Playbook**: Predefined playbook for creating users on remote hosts

## Prerequisites

- Go 1.25.5 or later
- Kafka broker running (default: localhost:9092)
- SSH access to target hosts

## Installation

1. Install dependencies:
```bash
go mod download
```

2. Build the application:
```bash
go build -o ansible-go
```

## Usage

### Start the Service

```bash
./ansible-go -brokers localhost:9092 -group ansible-go-group -topic ansible-tasks
```

Or with custom parameters:
```bash
./ansible-go -brokers kafka1:9092,kafka2:9092 -group my-group -topic my-tasks
```

### Send Tasks via Kafka

#### Example 1: Create User Playbook Format

Send a JSON message to the Kafka topic:

```json
{
  "ip": "10.1.1.1",
  "username": "user10",
  "user_id": "10",
  "ssh_user": "root",
  "ssh_pass": "your-password"
}
```

#### Example 2: Generic Task Format

```json
{
  "ip": "10.1.1.1",
  "action": "create_user",
  "params": {
    "username": "user10",
    "user_id": "10"
  },
  "username": "root",
  "password": "your-password"
}
```

### Using Kafka Producer (Example)

You can use any Kafka producer to send messages. Here's an example using `kafka-console-producer`:

```bash
kafka-console-producer --broker-list localhost:9092 --topic ansible-tasks
```

Then paste the JSON message and press Enter.

## Configuration

The application accepts the following command-line flags:

- `-brokers`: Kafka broker addresses (default: "localhost:9092")
- `-group`: Kafka consumer group ID (default: "ansible-go-group")
- `-topic`: Kafka topic to consume from (default: "ansible-tasks")

## Security Notes

- **SSH Keys**: For production, use SSH key-based authentication instead of passwords
- **Kafka Security**: Configure Kafka with SASL/SSL for production environments
- **Network**: Ensure proper firewall rules for SSH and Kafka connections

## Example Scenario

To create user "user10" with ID "10" on host "10.1.1.1":

1. Start the service:
```bash
./ansible-go
```

2. Send message to Kafka:
```json
{
  "ip": "10.1.1.1",
  "username": "user10",
  "user_id": "10",
  "ssh_user": "root",
  "ssh_pass": "your-password"
}
```

The service will:
- Connect to 10.1.1.1 via SSH
- Check if user "user10" exists
- Create the user with ID 10
- Set a default password

## Architecture

- `main.go`: Application entry point and Kafka listener setup
- `kafka_listener.go`: Kafka message consumption
- `handler.go`: Playbook and task processing
- `playbook.go`: Playbook definitions and parsing
- `executor.go`: SSH execution and remote command handling

## Future Enhancements

- Support for more playbook types
- Parallel task execution
- Task result reporting back to Kafka
- Playbook templates and variables
- Inventory management
- Task retry and error handling
