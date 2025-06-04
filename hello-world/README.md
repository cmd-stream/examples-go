# hello-world
hello-world example demonstrates the basic usage of cmd-stream-go. It includes 
definitions for:
- [SayHelloCmd](hello-world/cmds/say_hello.go) and [SayFancyHelloCmd](hello-world/cmds/say_fancy_hello.go) Commands.
- [Greeting](hello-world/results/greeting.go) Result.
- [Greeter](hello-world/receiver/greeter.go) Receiver.

## Details
- [musgen-go](github.com/mus-format/musgen-go) generates serializers for 
  `core.Cmd` and `core.Result` interfaces (see [cmds/gen](hello-world/cmds/gen) 
	and [results/gen](hello-world/results/gen)), to build client and 
	server codecs.
- Commands are sent to the server using the [sender](github.com/cmd-stream/sender-go), 
  which is built on top of the client group.
- Commands and Results are transmitted with deadlines.
- Sender is configured with [UnexpectedResultCallback](github.com/cmd-stream/core-go/client/client.go).
- Server is configured with [LostConnCallback](github.com/cmd-stream/core-go/server/server.go).
