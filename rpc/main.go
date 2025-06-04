package main

import (
	"context"
	"sync"

	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9005"

	// Start server.
	var (
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		codec   = cdc.NewServerCodec(results.ResultMUS, cmds.CmdMUS)
		wgS     = &sync.WaitGroup{}
	)
	server, err := utils.StartServer(addr, codec, greeter, wgS)
	assert.EqualError(err, nil)

	SayHello(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.Equal(err, nil)
}

func SayHello(addr string) {
	// Create sender.
	var (
		clientsCount = 1
		codec        = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	)
	sender, err := utils.MakeSender(addr, clientsCount, codec)
	assert.EqualError(err, nil)

	// Create service.
	service := GreeterService{sender}

	// Call SayHello.
	str, err := service.SayHello(context.Background(), "world")
	assert.EqualError(err, nil)
	assert.Equal(str, "Hello world")

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)
}
