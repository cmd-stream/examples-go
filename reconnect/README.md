# Reconnect Example

This example demonstrates how the `cmd-stream` library can automatically restore connections to the server using the `grp.WithReconnect` option. It shows a scenario where the server is intentionally closed and then restarted, and how the sender handles the reconnection.

## Features

- **Automatic Reconnection**: Shows how to configure a Sender with a client group that automatically attempts to reconnect upon connection loss.
- **Resilient Communication**: Demonstrates that once the server is back online, the Sender can resume sending Commands without manual intervention.

## Details

A connection can be lost while sending a Command or while waiting for a Result. In both cases, the reconnect mechanism allows the Command to be resent (assuming it is idempotent) after a short delay or upon the next attempt.

To enable reconnection in a Sender:
```go
sender, err := cmdstream.NewSender(addr, codec,
    sndr.WithGroup(
        grp.WithReconnect[T](),
    ),
)
```

## Running the Example

From the root of the `reconnect` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Closing server...
Starting server again...
Waiting for the sender to reconnect...
Sending "message"... Result: "message"
```
