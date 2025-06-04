# cmd-stream-examples-go
This repository contains several examples of using cmd-stream-go (each package 
is one example):
- [echo](examples-go/echo): A minimal example.
- [hello-world](examples-go/hello-world): Shows the basic usage of cmd-stream-go.
- [hello-world_protobuf](examples-go/hello-world_protobuf): Demonstrates the 
  basic usage of cmd-stream-go with the Protobuf serializer.
- [keepalive](examples-go/keepalive): Shows how the client can keep a connection 
  alive when there are no Commands to send.
- [reconnect](examples-go/reconnect): Demonstrates how the client can reconnect 
  to the server after losing the connection.
- [otel](examples-go/otel): Demonstrates how cmd-stream-go can be used with a 
  circuit breaker and OpenTelemetry.
- [server-streaming](examples-go/server-streaming): An example where the Command 
  sends back multiple Results.
- [client-group](examples-go/client-group): Demonstrates the usage of client 
  groups for high-performance communication with the server.
- [rpc](examples-go/rpc): Demonstrates how to implement RPC using cmd-stream-go.
- [tls](examples-go/tls): Demonstrates how to use cmd-stream-go with TLS.
- 