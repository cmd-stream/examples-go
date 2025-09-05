package main

import (
	"crypto/tls"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	utils "github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/handler-go"
	sndr "github.com/cmd-stream/sender-go"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	csrv "github.com/cmd-stream/core-go/server"
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
		greeter       = receiver.NewGreeter("Hello", "incredible", " ")
		invoker       = srv.NewInvoker(greeter)
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
	server := cmdstream.MakeServer(serverCodec, invoker,
		srv.WithCore(
			csrv.WithTLSConfig(&serverTLSConf),
		),
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
	sender, err := sndr.Make(addr, clientCodec,
		sndr.WithTLSConfig[receiver.Greeter](&clientTLSConf),
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

func SendCmd(sender sndr.Sender[receiver.Greeter]) {
	var (
		cmd  = cmds.NewSayHelloCmd("world")
		want = results.Greeting("Hello world")
	)
	greeting, err := utils.SendCmd(cmd, sender)
	assert.EqualError(err, nil)
	assert.Equal(greeting, want)
}
