package main

import (
	"sync"

	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/hello-world_protobuf/cmds"
	"github.com/cmd-stream/examples-go/hello-world_protobuf/results"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9002"

	// Start server.
	var (
		codec   = cdc.NewServerCodec(cmds.CmdProtobuf, results.ResultProtobuf)
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		wgS     = &sync.WaitGroup{}
	)
	server, err := utils.StartServer(addr, codec, greeter, wgS)
	assert.EqualError(err, nil)

	SendCmds(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendCmds(addr string) {
	// Create sender.
	codec := cdc.NewClientCodec(cmds.CmdProtobuf, results.ResultProtobuf)
	sender, err := utils.MakeSender(addr, 2, codec)
	assert.EqualError(err, nil)

	wgR := &sync.WaitGroup{}
	// Send SayHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = cmds.NewSayHelloCmd("world")
			wantGreeting = results.NewGreeting("Hello world")
		)
		err := utils.Exchange(sayHelloCmd, wantGreeting, sender)
		assert.EqualError(err, nil)
	}()

	// Send SayFancyHelloCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = cmds.NewSayFancyHelloCmd("world")
			wantGreeting     = results.NewGreeting("Hello incredible world")
		)
		err := utils.Exchange(sayFancyHelloCmd, wantGreeting, sender)
		assert.EqualError(err, nil)
	}()
	wgR.Wait()

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)
}
