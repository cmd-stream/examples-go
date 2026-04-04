# TLS Example

This example demonstrates how to use the `cmd-stream` library with the TLS protocol for secure communication between the server and the client.

## Features

- **TLS Configuration**: Shows how to load and use X.509 key pairs for both the server and the client.
- **Secure Communication**: Uses the `srv.WithCore` and `sndr.WithTLSConfig` options to enable TLS.

## Prerequisites

The example requires TLS certificates in the `certs/` directory:
- `server.pem` / `server.key`: Server's certificate and private key.
- `client.pem` / `client.key`: Client's certificate and private key.

## Running the Example

From the root of the `tls` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting TLS server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending "SayHelloCmd" with "world"... Result: "Hello world"
```
