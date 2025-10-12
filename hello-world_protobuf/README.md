# hello-world_protobuf

This example demonstrates how to use `cmd-stream-go` with Protocol Buffers.

## Differences From the hello-world Example

1. Commands (`SayHelloCmd` and `SayFancyHelloCmd`) store all properties in data
   structures that are serializable by Protobuf.
2. The `protobuf-format.go` files instead of `mus-format.gen.go`.

Everything else remains the same.
