# RPC Example

This example demonstrates how to implement a traditional RPC (Remote Procedure Call) interface on top of the `cmd-stream` library. It shows how to encapsulate Command sending logic within a service structure, providing a more familiar API for developers.

## Features

- **Service Abstraction**: Demonstrates wrapping a `sndr.Sender` within a [GreeterService](service.go) to hide the underlying Command-stream details.
- **Synchronous-like API**: Shows how a service method can wait for a Result and return it as a simple function return value.

## Key Components

- **[GreeterService](service.go)**: A service structure that provides a `SayHello` method, which internally sends a `SayHelloCmd` and returns the resulting `Greeting`.
- **[main.go](main.go)**: Shows how to initialize the server, the sender, and the service, followed by a simple service call.

## Running the Example

From the root of the `rpc` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Calling SayHello service...
Service responded with... Result: "Hello world"
```
