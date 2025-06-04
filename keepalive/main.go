package main

import (
	"sync"
	"time"

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
	const addr = "127.0.0.1:9003"

	// Start server.
	var (
		receiver = receiver.NewGreeter("Hello", "incredible", " ")
		codec    = cdc.NewServerCodec(results.ResultMUS, cmds.CmdMUS)
		wgS      = &sync.WaitGroup{}
	)
	server, err := utils.StartServer(addr, codec, receiver, wgS)
	assert.EqualError(err, nil)

	SendCmds(addr)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

func SendCmds(addr string) {
	// Create keepalive sender.
	var (
		codec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgR   = &sync.WaitGroup{}
	)
	sender, err := MakeKeepaliveSender(addr, 1, codec)
	assert.EqualError(err, nil)

	// Send SayHelloCmd command.
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

	// Ping-Pong time... When there are no Commands to send, clients of the sender
	// send a predefined PingCmd.
	time.Sleep(2 * utils.CmdReceiveDuration)

	// Send a command again after the long delay to check if the connection is
	// still alive.
	// Send SayFancyHelloCmd command.
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
