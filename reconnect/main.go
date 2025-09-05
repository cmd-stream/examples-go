package main

import (
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/handler-go"
	sndr "github.com/cmd-stream/sender-go"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	var (
		greeter     = receiver.NewGreeter("Hello", "incredible", " ")
		invoker     = srv.NewInvoker(greeter)
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)

	// Make server.
	server := cmdstream.MakeServer(serverCodec, invoker,
		srv.WithHandler(
			handler.WithAt(),
		),
	)
	// Start server.
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Make reconnect sender.
	sender, err := MakeReconnectSender(addr, clientCodec)
	assert.EqualError(err, nil)

	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)

	// Start the server again arter some time.
	time.Sleep(time.Second)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Wait for the sender clients to reconnect.
	time.Sleep(200 * time.Millisecond)

	// Send Command.
	SendCmd(sender)

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func MakeReconnectSender(addr string, codec cln.Codec[receiver.Greeter]) (
	sender sndr.Sender[receiver.Greeter], err error,
) {
	return sndr.Make(addr, codec,
		sndr.WithGroup[receiver.Greeter](
			grp.WithReconnect[receiver.Greeter](),
		),
	)
}

func SendCmd(sender sndr.Sender[receiver.Greeter]) {
	var (
		cmd  = cmds.NewSayHelloCmd("world")
		want = results.Greeting("Hello world")
	)
	greeting, err := utils.SendCmd(cmd, sender)
	assert.EqualError(err, nil)
	assert.Equal(greeting, want)
}
