package main

import (
	"context"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/handler-go"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	sndr "github.com/cmd-stream/sender-go"
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

	// Make sender.
	sender, err := sndr.Make(addr, clientCodec)
	assert.EqualError(err, nil)

	// Create service.
	service := GreeterService{sender}
	// Use service.
	str, err := service.SayHello(context.Background(), "world")
	assert.EqualError(err, nil)
	assert.Equal(str, "Hello world")

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}
