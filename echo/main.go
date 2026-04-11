package main

import (
	"fmt"
	"net"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	examples "github.com/cmd-stream/examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		serverCodec = ServerCodec{}
		clientCodec = ClientCodec{}
	)

	// Start server.
	fmt.Printf("Starting server on %s...\n", addr)
	go func() {
		server, _ := cmdstream.NewServer(struct{}{}, serverCodec)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make client.
	fmt.Println("Connecting to server...")
	client, err := MakeClient(addr, clientCodec)
	assert.EqualError(err, nil)

	// Send Command.
	SendCmd(client)

	// Close client.
	err = client.Close()
	assert.EqualError(err, nil)
}

func MakeClient(addr string, codec cln.Codec[struct{}]) (
	client *ccln.Client[struct{}], err error,
) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.NewClient(codec, conn)
}

func SendCmd(client *ccln.Client[struct{}]) {
	var (
		cmd          = examples.Message("one two three")
		asyncResults = make(chan core.AsyncResult, 1)
	)
	fmt.Printf("Sending message: \"%v\"\n", cmd)
	_, _, err := client.Send(cmd, asyncResults)
	assert.EqualError(err, nil)

	result := (<-asyncResults).Result.(examples.Message)
	fmt.Printf("Received echo... Result: %q\n", result)
	assert.Equal(cmd, result)
}
