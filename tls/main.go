package main

import (
	"sync"

	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	utils "github.com/cmd-stream/examples-go/hello-world/utils"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9007"

	// Start server.
	var (
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		codec   = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		wgS     = &sync.WaitGroup{}
	)
	server, err := StartServer(addr, codec, greeter, wgS)
	assert.EqualError(err, nil)

	SendCmd(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendCmd(addr string) {
	// Create sender.
	codec := cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	sender, err := MakeSender(addr, 1, codec)
	assert.EqualError(err, nil)

	// Send SayHelloCmd.
	wgR := &sync.WaitGroup{}
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
	wgR.Wait()

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)
}
