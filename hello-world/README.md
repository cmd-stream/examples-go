# hello-world
hello-world example demonstrates the basic usage of cmd-stream-go. It includes 
definitions for:
- [SayHelloCmd](examples-go/hello-world/cmds/say_hello.go) and [SayFancyHelloCmd](examples-go/hello-world/cmds/say_fancy_hello.go) Commands.
- [Greeting](examples-go/hello-world/results/greeting.go) Result.
- [Greeter](examples-go/hello-world/receiver/greeter.go) Receiver.

## Details
- [musgen-go](github.com/mus-format/musgen-go) generates serializers for 
  `core.Cmd` and `core.Result` interfaces (see [cmds/gen](examples-go/hello-world/cmds/gen) 
	and [results/gen](examples-go/hello-world/results/gen)), to build client and 
	server codecs.
- Commands are sent to the server using the [sender](github.com/cmd-stream/sender-go), 
  which is built on top of the client group.
- Commands and Results are transmitted with deadlines.
- Sender is configured with [UnexpectedResultCallback](github.com/cmd-stream/core-go/client/client.go).
- Server is configured with [LostConnCallback](github.com/cmd-stream/core-go/server/server.go).
