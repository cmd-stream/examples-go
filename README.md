# cmd-stream-examples-go
This repository contains several examples of using cmd-stream-go (each package 
is one example):
- [echo](echo): A minimal example.
- [hello-world](hello-world): Shows the basic usage of cmd-stream-go.
- [hello-world_protobuf](hello-world_protobuf): Demonstrates the 
  basic usage of cmd-stream-go with the Protobuf serializer.
- [keepalive](keepalive): Shows how the client can keep a connection 
  alive when there are no Commands to send.
- [reconnect](reconnect): Demonstrates how the client can reconnect 
  to the server after losing the connection.
- [otel](otel): Demonstrates how cmd-stream-go can be used with a 
  circuit breaker and OpenTelemetry.
- [server-streaming](server-streaming): An example where the Command 
  sends back multiple Results.
- [rpc](rpc): Demonstrates how to implement RPC using cmd-stream-go.
- [tls](tls): Demonstrates how to use cmd-stream-go with TLS.
- 