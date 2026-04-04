# Hello World Example

This is a fundamental example that demonstrates the basic usage of the `cmd-stream` library. It covers the core concepts: Commands, Results, and Receivers.

## Features

- **Basic Communication**: Shows how to send simple Commands and receive Results.
- **Custom Codecs**: Demonstrates the use of the MUS codec for efficient serialization.
- **Sender API**: Illustrates how to use the high-level `sndr.Sender` for Command delivery.

## Key Components

- **[SayHelloCmd](cmds/say_hello_cmd.go)** & **[SayFancyHelloCmd](cmds/say_fancy_hello_cmd.go)**: The Command definitions.
- **[Greeting](results/greeting.go)**: The Result type returned by the server.
- **[Greeter](receiver/greeter.go)**: The Receiver responsible for executing the Commands.
- **[codec-mus-stream-go](https://github.com/cmd-stream/codec-mus-stream-go)**: The MUS codec used for efficient Command and Result serialization.

## Running the Example

From the root of the `hello-world` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending "SayFancyHelloCmd" with "world"... Result: "Hello incredible world"
Sending "SayHelloCmd" with "world"... Result: "Hello world"
```
