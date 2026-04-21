package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	"github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdcjson "github.com/cmd-stream/codec-json-go"
	examples "github.com/cmd-stream/examples-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	serverCert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	assert.EqualError(err, nil)
	clientCert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	assert.EqualError(err, nil)

	var (
		registry = cdcjson.NewRegistry(
			cdcjson.WithCmd[struct{}, examples.Message](),
			cdcjson.WithResult[struct{}, examples.Message](),
		)
		serverCodec = cdcjson.NewServerCodecWith(registry)
		clientCodec = cdcjson.NewClientCodecWith(registry)

		serverTLSConf = tls.Config{Certificates: []tls.Certificate{serverCert}}
		clientTLSConf = tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true,
		}
		wgS = &sync.WaitGroup{}
	)

	// Make server.
	server, err := cmdstream.NewServer(struct{}{}, serverCodec,
		srv.WithCore(
			csrv.WithTLSConfig(&serverTLSConf),
		),
		srv.WithHandler(
			handler.WithAt(),
		),
	)
	assert.EqualError(err, nil)
	// Start server.
	fmt.Printf("Starting TLS server on %s...\n", addr)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Make sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := cmdstream.NewSender(addr, clientCodec,
		sndr.WithTLSConfig[struct{}](&clientTLSConf),
	)
	assert.EqualError(err, nil)

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
