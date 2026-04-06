package main

import (
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
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

const (
	KeepaliveTime  = 200 * time.Millisecond
	KeepaliveIntvl = 200 * time.Millisecond
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
			handler.WithCmdReceiveDuration(utils.CmdReceiveDuration),
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

	// Make keepalive sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := MakeKeepaliveSender(addr, clientCodec)
	assert.EqualError(err, nil)
	// Send Commands.
	SendCmds(sender)

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func MakeKeepaliveSender(addr string, codec cln.Codec[rcvr.Greeter]) (
	sender sndr.Sender[rcvr.Greeter], err error,
) {
	return cmdstream.NewSender(addr, codec,
		sndr.WithGroup(
			grp.WithClient[rcvr.Greeter](
				cln.WithKeepalive(
					dcln.WithKeepaliveTime(KeepaliveTime),
					dcln.WithKeepaliveIntvl(KeepaliveIntvl),
				),
			),
		),
	)
}

func SendCmds(sender sndr.Sender[rcvr.Greeter]) {
	// Send SayHelloCmd.
	{
		var (
			cmd  = cmds.SayHelloCmd{Str: "world"}
			want = results.Greeting("Hello world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		fmt.Printf("Sending \"SayHelloCmd\" with \"world\"... Result: %q\n", greeting)
		assert.Equal(greeting, want)
	}

	// Ping-Pong time... When there are no Commands to send, clients will send
	// PingCmd to the server to maintain the connection.
	fmt.Println("Ping-Pong time...")
	time.Sleep(2 * utils.CmdReceiveDuration)

	// Still able to send SayFancyHelloCmd.
	{
		var (
			cmd  = cmds.SayFancyHelloCmd{Str: "world"}
			want = results.Greeting("Hello incredible world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		fmt.Printf("Sending \"SayFancyHelloCmd\" with \"world\"... Result: %q\n", greeting)
		assert.Equal(greeting, want)
	}
}
