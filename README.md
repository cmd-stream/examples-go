# `cmd-stream` Examples for Go

This repository contains a collection of examples demonstrating how to use the
[`cmd-stream`](https://github.com/cmd-stream/cmd-stream-go) library in Go. From
simple "Hello World" scenarios to complex observability and resilience setups,
these examples provide a comprehensive guide to building modular, command-based
systems.

## Getting Started

If you're new to `cmd-stream`, a good starting point is
**[`calc_json`](calc_json)** and **[`hello-world`](hello-world)**.

## Example Index

The examples are organized into logical categories to help you find what you
need:

### Basic Communication

- **[`echo`](echo)**: Minimal example showing direct client usage (without a
  `Sender`).
- **[`calc_json`](calc_json)**: Simple request-response pattern using the JSON
  codec.
- **[`calc_protobuf`](calc_protobuf)**: Efficient communication using Protobuf.
- **[`hello-world`](hello-world)**: Demonstrates core concepts with the MUS codec.

### Advanced Features

- **[`server-streaming`](server-streaming)**: Multi-Result Commands for
  server-side streaming responses.
- **[`tls`](tls)**: Secure communication between server and client using the TLS
  protocol.
- **[`keepalive`](keepalive)**: Maintaining idle connections via heartbeat
  (Ping) mechanisms.
- **[`reconnect`](reconnect)**: Automatic reconnection logic.

### Specialized Patterns

- **[`rpc`](rpc)**: Shows how to implement traditional RPC-style service
interfaces on top of `cmd-stream`.
- **[`otel`](otel)**: Full observability stack (Tracing, Metrics) using
  OpenTelemetry and a Circuit Breaker.

## Common Architecture

Most examples follow a standard pattern using the high-level `cmd-stream` API:

```go
// Server initialization
server, _ := cmdstream.NewServer(receiver, codec, ...)
server.ListenAndServe(addr)

// Client initialization (Sender)
sender, _ := cmdstream.NewSender(addr, codec, ...)
result, _ := sender.Send(context.Background(), myCommand)
```

## Running the Examples

Each example resides in its own directory and includes a dedicated `README.md`
with specific details. Generally, you can run any example by moving into its
folder and using `go run`:

```bash
cd [example-folder]
go run .
```
