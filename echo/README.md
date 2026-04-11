# Echo Example

This is a simple example that demonstrates the basic request-response pattern using the `cmd-stream` library. Unlike other examples that use the high-level `Sender`, this example shows how to use the low-level `Client` directly.

## Features

- **Direct Client Usage**: Demonstrates how to create and use a `ccln.Client` for sending individual Commands.
- **Async Results**: Shows how to handle Results via Go channels using the `client.Send` method.

## Key Components

- **[Message](../message.go)**: A simple string-based Command and Result type.
- **[ClientCodec](codecs.go)** & **[ServerCodec](codecs.go)**: Simple custom codecs for the `Message` type.

## Running the Example

From the root of the `echo` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Connecting to server...
Sending message: "one two three"
Received echo... Result: "one two three"
```
