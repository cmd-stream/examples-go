package main

import (
	"net"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	assert "github.com/ymz-ncnk/assert/panic"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		invoker     = srv.NewInvoker(struct{}{})
		serverCodec = ServerCodec{}
		clientCodec = ClientCodec{}
	)

	// Start server.
	go func() {
		server := cmdstream.MakeServer(serverCodec, invoker)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make client.
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
	return cmdstream.MakeClient(codec, conn)
}

func SendCmd(client *ccln.Client[struct{}]) {
	var (
		cmd          = Message("one two three")
		asyncResults = make(chan core.AsyncResult, 1)
	)
	_, _, err := client.Send(cmd, asyncResults)
	assert.EqualError(err, nil)

	result := (<-asyncResults).Result.(Message)
	assert.Equal(cmd, result)
}
