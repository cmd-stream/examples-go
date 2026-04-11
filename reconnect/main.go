package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-json-go"
	examples "github.com/cmd-stream/examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	var (
		cmdTypes = []reflect.Type{
			reflect.TypeFor[examples.Message](),
		}
		resultTypes = []reflect.Type{
			reflect.TypeFor[examples.Message](),
		}
		serverCodec = cdc.NewServerCodec[struct{}](cmdTypes, resultTypes)
		clientCodec = cdc.NewClientCodec[struct{}](cmdTypes, resultTypes)
		wgS         = &sync.WaitGroup{}
	)

	// Make server.
	server, err := cmdstream.NewServer(struct{}{}, serverCodec,
		srv.WithHandler(
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

	// Make reconnect sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := MakeReconnectSender(addr, clientCodec)
	assert.EqualError(err, nil)

	// Close server.
	fmt.Println("-- Closing server... --")
	err = server.Close()
	assert.EqualError(err, nil)

	// Start the server again after some time.
	time.Sleep(time.Second)
	fmt.Println("-- Starting server again... --")
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

func MakeReconnectSender(addr string, codec cln.Codec[struct{}]) (
	sender sndr.Sender[struct{}], err error,
) {
	return cmdstream.NewSender(addr, codec,
		sndr.WithGroup(
			grp.WithReconnect[struct{}](),
		),
	)
}

func SendCmd(sender sndr.Sender[struct{}]) {
	var (
		cmd  = examples.Message("message")
		want = examples.Message("message")
	)
	result, err := sender.Send(context.Background(), cmd)
	assert.EqualError(err, nil)
	fmt.Printf("Sending \"%v\"... Result: \"%v\"\n", cmd, result)
	assert.Equal(result.(examples.Message), want)
}
