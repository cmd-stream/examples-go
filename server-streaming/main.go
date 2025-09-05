package main

import (
	"context"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"
	"github.com/cmd-stream/examples-go/server-streaming/results"
	"github.com/cmd-stream/handler-go"
	sndr "github.com/cmd-stream/sender-go"

	srv "github.com/cmd-stream/cmd-stream-go/server"
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

	// Make keepalive sender.
	sender, err := sndr.Make(addr, clientCodec)
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

func SendMultiCmd(sender sndr.Sender[receiver.Greeter]) {
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
			assert.Equal(greeting, want[i])
			i++
			return nil
		}
		deadline = time.Now().Add(utils.CmdSendDuration)
	)
	defer cancel()
	err := sender.SendMultiWithDeadline(ctx, cmd, len(want), resultHandler,
		deadline)
	assert.EqualError(err, nil)
}
