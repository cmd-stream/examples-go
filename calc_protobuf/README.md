# Protobuf Calculator Example

This example demonstrates how to use the `cmd-stream` library with **Protobuf** 
for efficient and type-safe serialization. It implements a simple calculator 
that can perform addition and subtraction.

## Key Components

- **[AddCmd](cmds/add_cmd.go)** & **[SubCmd](cmds/sub_cmd.go)**: Provide `core.Cmd` 
  interface implementations for the Protobuf-generated messages.
- **[CalcResult](results/calc_result.go)**: Provides `core.Result` interface 
  implementation for the Protobuf-generated message.
- **[Calc](receiver/calc.go)**: The Receiver responsible for performing the 
  mathematical operations.
- **[cmds.proto](cmds/cmds.proto)** & **[results.proto](results/results.proto)**: The Protobuf definition files.
- **[codec-protobuf-go](https://github.com/cmd-stream/codec-protobuf-go)**: The Protobuf codec used for Command and Result serialization.

## Running the Example

From the root of the `calc_protobuf` directory, run:
```bash
go run .
```

## Expected Output

```text
Starting server on 127.0.0.1:9000...
Initializing sender and connecting...
Sending AddCmd(2, 3)... Result: 5
Sending SubCmd(8, 4)... Result: 4
```
