package main

import (
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/handler-go"
	sndr "github.com/cmd-stream/sender-go"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
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
		invoker     = srv.NewInvoker[rcvr.Greeter](greeter)
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)

	// Make server.
	server := cmdstream.MakeServer(serverCodec, invoker,
		srv.WithHandler(
			handler.WithCmdReceiveDuration(utils.CmdReceiveDuration),
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

	// Make keepalive sender.
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
	return sndr.Make(addr, codec,
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
	wg := sync.WaitGroup{}

	// Send SayHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.NewSayHelloCmd("world")
			want = results.Greeting("Hello world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting, want)
		wg.Done()
	}()

	// Ping-Pong time... When there are no Commands to send, clients will send
	// PingCmd.
	time.Sleep(2 * utils.CmdReceiveDuration)

	// Send SayFancyHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.NewSayFancyHelloCmd("world")
			want = results.Greeting("Hello incredible world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting, want)
		wg.Done()
	}()

	wg.Wait()
}
