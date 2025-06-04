# otel
This example demonstrates how cmd-stream-go can be used with a circuit breaker 
and OpenTelemetry. It runs for 90 seconds, during which the client sends both 
regular and traceable Commands to the server. The server is available only 
during the first and last 30 seconds. During the middle 30 seconds, it is down, 
allowing you to observe the circuit breaker in action through log messages.

```
                         circuit breaker active
|-----------------------|-----------------------|-----------------------|
        server up              server down              server up 

```

Spans are created for each Command on the client side and propagated to the 
server using `TraceCmd`, a simple wrapper around regular Commands.

## How It Works
- The app sends traces and metrics to the OTEL Collector via OTLP (port 4317).
- The OTEL Collector exports traces to Tempo (OTLP on tempo:4317).
- The OTEL Collector exposes metrics on port 8889 (Prometheus format).
- Prometheus scrapes metrics from the OTEL Collector (otel-collector:8889).
- Grafana reads traces from Tempo and metrics from Prometheus.

## How to Run
Make sure:
1. Docker is installed and running on your system.
2. Ports 9090, 3200, 4317, and 3000 are available.
   
From the `otel/monitoring` folder, run:
```bash
$ docker compose up -d
```

From the `otel` folder, run:
```bash
$ go run .
```

### Access
- Grafana UI: http://localhost:3000 (admin/admin)
- Prometheus UI: http://localhost:9090
