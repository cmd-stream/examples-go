# Keepalive Example

This example demonstrates how the `cmd-stream` library can keep a connection alive even when there are no Commands being sent. This is achieved by configuring the clients to send periodic heartbeat (Ping) Commands.

## Features

- **Heartbeat Mechanism**: Shows how to use the `cln.WithKeepalive` option to prevent connection timeouts during idle periods.
- **Inactivity Detection**: Demonstrates how the client can automatically detect inactivity and start sending Ping Commands.

## Details

To enable keepalive in a Sender:
```go
sender, err := cmdstream.NewSender(addr, codec,
    sndr.WithGroup(
        grp.WithClient[T](
            cln.WithKeepalive(
                dcln.WithKeepaliveTime(KeepaliveTime),
                dcln.WithKeepaliveIntvl(KeepaliveIntvl),
            ),
        ),
    ),
)
```

## Running the Example

From the root of the `keepalive` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending "SayHelloCmd" with "world"... Result: "Hello world"
Ping-Pong time...
Sending "SayFancyHelloCmd" with "world"... Result: "Hello incredible world"
```
