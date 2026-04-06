package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"
	"github.com/cmd-stream/examples-go/server-streaming/results"
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

	// Make keepalive sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := cmdstream.NewSender(addr, clientCodec)
	assert.EqualError(err, nil)
	// Send Command that sends back several Results.
	SendMultiCmd(sender)

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func SendMultiCmd(sender sndr.Sender[rcvr.Greeter]) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(),
			utils.ResultReceiveDuration)
		cmd  = cmds.NewSayFancyHelloMultiCmd("world")
		want = []results.Greeting{
			results.NewGreeting("Hello", false),
			results.NewGreeting("incredible", false),
			results.NewGreeting("world", true),
		}
		i             int
		resultHandler sndr.ResultHandlerFn = func(result core.Result, err error) error {
			greeting := result.(results.Greeting)
			fmt.Printf("Received... Result: %q\n", greeting)
			assert.Equal(greeting, want[i])
			i++
			return nil
		}
		deadline = time.Now().Add(utils.CmdSendDuration)
	)
	defer cancel()
	fmt.Println("Sending multi-result command...")
	err := sender.SendMultiWithDeadline(ctx, deadline, cmd, len(want), resultHandler)
	assert.EqualError(err, nil)
}
