# JSON Calculator Example (Getting Started)

This is a simple "getting started" example that demonstrates how to implement a basic request-response pattern using **JSON** for serialization. It implements a calculator that can perform addition and subtraction.

## Key Components

- **[AddCmd](cmds.go)** & **[SubCmd](cmds.go)**: The Command definitions.
- **[CalcResult](results.go)**: The Result type.
- **[Calc](receiver.go)**: The Receiver responsible for performing the mathematical operations.
- **[codec-json-go](https://github.com/cmd-stream/codec-json-go)**: The JSON codec used for Command and Result serialization.

## Running the Example

From the root of the `calc_json` directory, run:
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
