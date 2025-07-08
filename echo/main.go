package main

import (
	"net"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	csrv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/echo/cmds"
	"github.com/cmd-stream/examples-go/echo/results"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"

	// Start server.
	l, err := net.Listen("tcp", addr)
	assert.EqualError(err, nil)
	server := cmdstream.MakeServer(ServerCodec{}, csrv.NewInvoker(struct{}{}))
	go func() {
		server.Serve(l.(*net.TCPListener))
	}()

	// Create client.
	conn, err := net.Dial("tcp", addr)
	assert.EqualError(err, nil)
	client, err := cmdstream.MakeClient(ClientCodec{}, conn)
	assert.EqualError(err, nil)

	// Send a Command and get the Result.
	var (
		str          = "Hello world"
		EchoCmd      = cmds.EchoCmd(str)
		wantResult   = results.EchoResult(str)
		asyncResults = make(chan core.AsyncResult, 1)
	)
	_, _, err = client.Send(EchoCmd, asyncResults)
	assert.EqualError(err, nil)
	assert.Equal((<-asyncResults).Result.(results.EchoResult), wantResult)

	// Close client.
	err = client.Close()
	assert.EqualError(err, nil)

	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
}
