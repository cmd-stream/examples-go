package main

import (
	"sync"

	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"
	"github.com/cmd-stream/examples-go/server-streaming/results"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9006"

	// Start server.
	var (
		receiver = receiver.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(results.ResultMUS, cmds.CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := utils.StartServer(addr, codec, receiver, wgS)
	assert.EqualError(err, nil)

	SendMultiCmd(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendMultiCmd(addr string) {
	// Create sender.
	codec := cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	sender, err := utils.MakeSender(addr, 1, codec)
	assert.EqualError(err, nil)

	wgR := &sync.WaitGroup{}
	// Send SayFancyHelloMultiCmd command.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayFancyHelloMultiCmd = cmds.NewSayFancyHelloMultiCmd("world")
			wantGreetings         = []results.Greeting{
				results.NewGreeting("Hello", false),
				results.NewGreeting("incredible", false),
				results.NewGreeting("world", true),
			}
		)
		err := Exchange(sayFancyHelloMultiCmd, wantGreetings, sender)
		assert.EqualError(err, nil)
	}()
	wgR.Wait()

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)
}
