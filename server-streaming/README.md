# Server Streaming Example

This example demonstrates how to use the `cmd-stream` library to send multiple Results back to the client in response to a single Command, often referred to as "server-side streaming".

## Features

- **Multi-Result Commands**: Shows how to define a Command that produces multiple Results.
- **Client-Side Handling**: Demonstrates the use of `sndr.ResultHandlerFn` to process a stream of Results asynchronously.

## Key Components

- **[SayFancyHelloMultiCmd](cmds/say_fancy_hello_multi_cmd.go)**: A Command that decomposes a string and sends each part back as a separate Result.
- **[Greeting](results/greeting.go)**: The Result type used for each streamed message.

## Running the Example

From the root of the `server-streaming` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending multi-result command...
Received... Result: "Hello"
Received... Result: "incredible"
Received... Result: "world"
```
