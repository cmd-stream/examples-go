# hello-world
hello-world example demonstrates the basic usage of cmd-stream-go. It includes 
definitions for:
- [SayHelloCmd](cmds/say_hello.go) and [SayFancyHelloCmd](cmds/say_fancy_hello.go) Commands.
- [Greeting](results/greeting.go) Result.
- [Greeter](receiver/greeter.go) Receiver.

## Details

- [musgen-go](https://github.com/mus-format/musgen-go) generates serializers for
  `core.Cmd` and `core.Result` interfaces (see [cmds/gen](cmds/gen) 
	and [results/gen](results/gen)), to build client and server codecs.
- Commands are sent to the server using the [sender](https://github.com/cmd-stream/sender-go), 
  which is built on top of the client group.
- Commands and Results are transmitted with deadlines.
- Sender is configured with [UnexpectedResultCallback](https://github.com/cmd-stream/core-go/client/client.go).
- Server is configured with [LostConnCallback](https://github.com/cmd-stream/core-go/server/server.go).
