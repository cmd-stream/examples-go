package main

import (
	"sync"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"

	// Start server.
	var (
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		// Serializers for core.Cmd and core.Result interfaces allow building
		// a server codec.
		codec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		wgS   = &sync.WaitGroup{}
	)

	server, err := utils.StartServer(addr, codec, greeter, wgS)
	assert.EqualError(err, nil)

	SendCmds(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

// SendCmds sends two Commands concurrently.
func SendCmds(addr string) {
	// Create sender.
	var (
		codec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgR   = &sync.WaitGroup{}
	)
	sender, err := utils.MakeSender(addr, 2, codec)
	assert.EqualError(err, nil)

	// Send SayHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = cmds.NewSayHelloCmd("world")
			wantGreeting = results.Greeting("Hello world")
		)
		err := utils.Exchange(sayHelloCmd, wantGreeting, sender)
		assert.EqualError(err, nil)
	}()

	// Send SayFancyHelloCmd.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloCmd = cmds.NewSayFancyHelloCmd("world")
			wantGreeting     = results.Greeting("Hello incredible world")
		)
		err := utils.Exchange(sayFancyHelloCmd, wantGreeting, sender)
		assert.EqualError(err, nil)
	}()
	wgR.Wait()

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)
}
