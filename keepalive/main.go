package main

import (
	"context"
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
	cdcjson "github.com/cmd-stream/codec-json-go"
	examples "github.com/cmd-stream/examples-go"
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
		registry = cdcjson.NewRegistry(
			cdcjson.WithCmd[struct{}, examples.Message](),
			cdcjson.WithResult[struct{}, examples.Message](),
		)
		serverCodec = cdcjson.NewServerCodecWith(registry)
		clientCodec = cdcjson.NewClientCodecWith(registry)

		wgS = &sync.WaitGroup{}
	)

	// Make server.
	server, err := cmdstream.NewServer(struct{}{}, serverCodec,
		srv.WithHandler(
			handler.WithCmdReceiveDuration(500*time.Millisecond),
			handler.WithAt(),
		),
	)
	assert.EqualError(err, nil)
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

func MakeKeepaliveSender(addr string, codec cln.Codec[struct{}]) (
	sender sndr.Sender[struct{}], err error,
) {
	return cmdstream.NewSender(addr, codec,
		sndr.WithGroup(
			grp.WithClient[struct{}](
				cln.WithKeepalive(
					dcln.WithKeepaliveTime(KeepaliveTime),
					dcln.WithKeepaliveIntvl(KeepaliveIntvl),
				),
			),
		),
	)
}

func SendCmds(sender sndr.Sender[struct{}]) {
	// Send message.
	{
		var (
			cmd  = examples.Message("one")
			want = examples.Message("one")
		)
		result, err := sender.Send(context.Background(), cmd)
		assert.EqualError(err, nil)
		fmt.Printf("Sending \"%v\"... Result: \"%v\"\n", cmd, result)
		assert.Equal(result.(examples.Message), want)
	}

	// Ping-Pong time... When there are no Commands to send, clients will send
	// PingCmd to the server to maintain the connection.
	fmt.Println("Ping-Pong time...")
	time.Sleep(1 * time.Second)

	// Still able to send message.
	{
		var (
			cmd  = examples.Message("two")
			want = examples.Message("two")
		)
		result, err := sender.Send(context.Background(), cmd)
		assert.EqualError(err, nil)
		fmt.Printf("Sending \"%v\"... Result: \"%v\"\n", cmd, result)
		assert.Equal(result.(examples.Message), want)
	}
}
