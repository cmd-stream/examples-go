# OpenTelemetry & Circuit Breaker Example

This example demonstrates how to use the `cmd-stream` library with 
**OpenTelemetry (OTel)** and a **Circuit Breaker** for observability and 
resilience.

## Features

- **Distributed Tracing**: Shows how to propagate trace context from client to 
  server using `otelcmd.TraceCmd`.
- **Circuit Breaking**: Demonstrates using the [circbrk](https://github.com/ymz-ncnk/circbrk-go) 
  library with Sender hooks to prevent overloading a failing server.
- **Server Restart Scenario**: Animates a realistic scenario where the server 
  goes down and comes back up, allowing you to observe the Circuit Breaker state 
  transitions.

```text
                          circuit breaker open
|-----------------------|-----------------------|-----------------------|
        server up              server down              server up 
         (30s)                  (30s)                    (30s)
```

## How It Works

1. **Distributed Tracing**: Spans are created for each Command on the client and propagated to the server.
2. **Infrastructure**:
    - **OTEL Collector**: Receives traces and metrics via OTLP (port 4317).
    - **Tempo**: Stores traces.
    - **Prometheus**: Scrapes and stores metrics.
    - **Grafana**: Visualizes traces from Tempo and metrics from Prometheus.

## Prerequisites

- **Docker**: Required to run the monitoring stack (Grafana, Tempo, Prometheus, 
  OTEL Collector).

## Running the Example

### 1. Start the Monitoring Stack

From the `otel/monitoring` folder, run:
```bash
docker compose up -d
```
*(If you encounter permission issues, try `sudo docker compose up -d` or the 
older `docker-compose` syntax)*

### 2. Run the Code

From the `otel` folder, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Processing commands with server restarts...
2026/04/04 13:21:12 -- close server --
2026/04/04 13:21:15 CircuitBreaker: Open
2026/04/04 13:21:21 CircuitBreaker: HalfOpen
...
2026/04/04 13:21:42 -- start server --
2026/04/04 13:21:45 CircuitBreaker: HalfOpen
2026/04/04 13:21:45 CircuitBreaker: Closed
```

## Accessing Visualizations

- **Grafana UI**: <http://localhost:3000> (Login: `admin`/`admin`)
- **Prometheus UI**: <http://localhost:9090>

## Shutdown

From the `otel/monitoring` folder:
```bash
docker compose down
```
