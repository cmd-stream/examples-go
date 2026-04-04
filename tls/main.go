package main

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/cmd-stream-go/handler"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	utils "github.com/cmd-stream/examples-go/hello-world/utils"

	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
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
		greeter       = rcvr.NewGreeter("Hello", "incredible", " ")
		serverCodec   = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec   = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		serverTLSConf = tls.Config{Certificates: []tls.Certificate{serverCert}}
		clientTLSConf = tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true,
		}
		wgS = &sync.WaitGroup{}
	)

	// Make server.
	server, _ := cmdstream.NewServer(greeter, serverCodec,
		srv.WithCore(
			csrv.WithTLSConfig(&serverTLSConf),
		),
		srv.WithHandler(
			handler.WithAt(),
		),
	)
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
		sndr.WithTLSConfig[rcvr.Greeter](&clientTLSConf),
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

func SendCmd(sender sndr.Sender[rcvr.Greeter]) {
	var (
		cmd  = cmds.SayHelloCmd{Str: "world"}
		want = results.Greeting("Hello world")
	)
	greeting, err := utils.SendCmd(cmd, sender)
	assert.EqualError(err, nil)
	fmt.Printf("Sending \"SayHelloCmd\" with \"world\"... Result: %q\n", greeting)
	assert.Equal(greeting, want)
}
