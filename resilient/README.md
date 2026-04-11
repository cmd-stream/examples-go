# Resilient Example

This example demonstrates how to build a highly available and resilient client 
using the `cmd-stream` library. It showcases the integration of multiple 
reliability patterns including automatic reconnection, keepalives, and circuit 
breaking.

## Features

- **High-Level Sender**: Uses the `sndr.Sender` abstraction which manages a 
  group of clients and high-level delivery logic.
- **Multiple Clients**: Configured with a group of 2 clients to provide parallel
  processing and redundancy.
- **Automatic Reconnection**: Leverages `grp.WithReconnect` to automatically 
  restore connections if the server goes down or the network is interrupted.
- **Active Keepalives**: Prevents the server from closing idle connections by 
  sending periodic heartbeats.
- **Circuit Breaker**: Integrates `circbrk-go` via `hks.CircuitBreakerHooks` to 
  stop sending traffic when a failure threshold is reached.

## How it Works

The example follows this sequence:
1.  Starts a Server and connects a resilient Sender.
2.  Sends a successful message to verify the connection is healthy.
3.  Closes the Server to simulate a failure.
4.  Attempts to send messages, which fail and trigger the Circuit Breaker 
    (transitioning to the `Open` state).
5.  Restarts the Server.
6.  Sends the message again (`message5`) to close the Circuit Breaker.

## Running the Example

From the root of the `resilient` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending message1... Result: message1
-- Closing server... --
Sending message2... Error: EOF
Sending message3... Error: timeout
Circuit breaker state changed to: Open
Sending message4... Error: timeout
-- Starting server on 127.0.0.1:9000... --
Circuit breaker state changed to: HalfOpen
Circuit breaker state changed to: Closed
Sending message5... Result: message5
```
