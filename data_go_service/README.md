# Data Go Service

This is a high-standard Go microservice designed to consume data from two RabbitMQ queues, process it, and forward it to HTTP endpoints.

## Project Structure

- `cmd/server/main.go`: Application entry point.
- `config/`: Configuration management using environment variables.
- `internal/rabbitmq/`: RabbitMQ connection and manual-acknowledgement consumer.
- `internal/processor/`: Business logic functions for each queue.
- `internal/httpclient/`: HTTP client for POSTing processed data.
- `internal/models/`: Shared JSON structures.
- `pkg/logger/`: Structured logging utility.

## Event-Driven Flow

1. **Message Arrival**: Triggered when a value enters `queue1` or `queue2`.
2. **Processing**: 
   - `queue1` messages are handled by `processor.Function1`.
   - `queue2` messages are handled by `processor.Function2`.
3. **HTTP Dispatch**: Processed data is sent to the corresponding HTTP endpoint (`http1` or `http2`).
4. **Acknowledgement**: The message is only acknowledged (`ack`) in RabbitMQ after a successful HTTP response.

## Configuration

Set the following environment variables (or use defaults in `config/config.go`):

- `RABBITMQ_URI`: Connection string (default: `amqp://guest:guest@localhost:5672/`)
- `QUEUE1_NAME`: Name of first queue (default: `queue1`)
- `QUEUE2_NAME`: Name of second queue (default: `queue2`)
- `HTTP1_URL`: URL for first queue's endpoint (default: `http://example.com/api1`)
- `HTTP2_URL`: URL for second queue's endpoint (default: `http://example.com/api2`)

## How to Run

1. Ensure Go is installed.
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the service:
   ```bash
   go run cmd/server/main.go
   ```
4. Run the mock server:
   ```bash
   go run cmd/mock_server/main.go
   ```
