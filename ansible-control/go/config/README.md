# Configuration

## Config File Location

The configuration file is located at:
```
ansible-control/go/config/config.json
```

## Configuration Format

Edit `config.json` to change Kafka settings:

```json
{
  "kafka": {
    "broker": "localhost:29092",
    "consumer_topic": "consumer-topic",
    "producer_topic": "producer-topic"
  }
}
```

## Configuration Options

### `kafka.broker`
- **Description**: Kafka broker address
- **Format**: `host:port`
- **Example**: `"localhost:29092"` or `"kafka.example.com:9092"`

### `kafka.consumer_topic`
- **Description**: Kafka topic to consume events from
- **Example**: `"consumer-topic"` or `"ansible-events"`

### `kafka.producer_topic`
- **Description**: Kafka topic to send responses to
- **Example**: `"producer-topic"` or `"ansible-responses"`

## How to Change Configuration

1. **Edit the config file:**
   ```bash
   # Edit config.json
   notepad config\config.json    # Windows
   # or
   nano config/config.json       # Linux/macOS
   ```

2. **Update the values:**
   ```json
   {
     "kafka": {
       "broker": "your-kafka-server:9092",
       "consumer_topic": "your-consumer-topic",
       "producer_topic": "your-producer-topic"
     }
   }
   ```

3. **Restart the agent:**
   - Stop the current agent (Ctrl+C)
   - Run it again using deploy scripts

## Example Configurations

### Local Development:
```json
{
  "kafka": {
    "broker": "localhost:29092",
    "consumer_topic": "consumer-topic",
    "producer_topic": "producer-topic"
  }
}
```

### Production:
```json
{
  "kafka": {
    "broker": "kafka.production.com:9092",
    "consumer_topic": "ansible-commands",
    "producer_topic": "ansible-results"
  }
}
```

### Multiple Brokers:
```json
{
  "kafka": {
    "broker": "kafka1:9092,kafka2:9092,kafka3:9092",
    "consumer_topic": "consumer-topic",
    "producer_topic": "producer-topic"
  }
}
```

## Notes

- The config file is loaded when the agent starts
- Changes require restarting the agent to take effect
- The config file path is automatically detected (checks multiple locations)
- If config file is not found, the application will fail to start
