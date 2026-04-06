package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/handler"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
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

	// Make sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := cmdstream.NewSender(addr, clientCodec)
	assert.EqualError(err, nil)

	// Create service.
	service := GreeterService{sender}
	// Use service.
	fmt.Println("Calling SayHello service...")
	str, err := service.SayHello(context.Background(), "world")
	assert.EqualError(err, nil)
	fmt.Printf("Service responded with... Result: %q\n", str)
	assert.Equal(str, "Hello world")

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}
