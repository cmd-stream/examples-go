package main

import (
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	var (
		greeter     = rcvr.NewGreeter("Hello", "incredible", " ")
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)

	// Make server.
	server, _ := cmdstream.NewServer(greeter, serverCodec,
		srv.WithHandler(
			handler.WithAt(),
		),
	)
	// Start server.
	fmt.Printf("Starting server on %s...\n", addr)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Make reconnect sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := MakeReconnectSender(addr, clientCodec)
	assert.EqualError(err, nil)

	// Close server.
	fmt.Println("Closing server...")
	err = server.Close()
	assert.EqualError(err, nil)

	// Start the server again arter some time.
	time.Sleep(time.Second)
	fmt.Println("Starting server again...")
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Wait for the sender clients to reconnect.
	fmt.Println("Waiting for the sender to reconnect...")
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

func MakeReconnectSender(addr string, codec cln.Codec[rcvr.Greeter]) (
	sender sndr.Sender[rcvr.Greeter], err error,
) {
	return cmdstream.NewSender(addr, codec,
		sndr.WithGroup(
			grp.WithReconnect[rcvr.Greeter](),
		),
	)
}

func SendCmd(sender sndr.Sender[rcvr.Greeter]) {
	var (
		cmd  = cmds.SayHelloCmd{Str: "world"}
		want = results.Greeting("Hello world")
	)
	greeting, err := utils.SendCmd(cmd, sender)
	assert.EqualError(err, nil)
	fmt.Printf("Sending \"SayHelloCmd\" with \"world\"... Result: %q\n", greeting)
	assert.Equal(greeting, want)
}
