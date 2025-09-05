package main

import (
	"context"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	utils "github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/hello-world_protobuf/cmds"
	"github.com/cmd-stream/examples-go/hello-world_protobuf/results"
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
		serverCodec = cdc.NewServerCodec(cmds.CmdProtobuf, results.ResultProtobuf)
		clientCodec = cdc.NewClientCodec(cmds.CmdProtobuf, results.ResultProtobuf)
		wgS         = &sync.WaitGroup{}
	)

	server := cmdstream.MakeServer(serverCodec, invoker,
		srv.WithHandler(
			handler.WithAt(),
		),
	)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	sender, err := sndr.Make(addr, clientCodec)
	assert.EqualError(err, nil)
	SendCmds(sender)

	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func SendCmds(sender sndr.Sender[receiver.Greeter]) {
	wg := sync.WaitGroup{}

	// Send SayHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.NewSayHelloCmd("world")
			want = results.NewGreeting("Hello world")
		)
		greeting, err := SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting.String(), want.String())
		wg.Done()
	}()

	// Send SayFancyHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.NewSayFancyHelloCmd("world")
			want = results.NewGreeting("Hello incredible world")
		)
		greeting, err := SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting.String(), want.String())
		wg.Done()
	}()

	wg.Wait()
}

func SendCmd[T any](cmd core.Cmd[T], sender sndr.Sender[T]) (
	greeting results.Greeting, err error,
) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(),
			utils.ResultReceiveDuration)
		deadline = time.Now().Add(utils.CmdSendDuration)
	)
	defer cancel()
	result, err := sender.SendWithDeadline(ctx, cmd, deadline)
	if err != nil {
		return
	}
	greeting = result.(results.Greeting)
	return
}
