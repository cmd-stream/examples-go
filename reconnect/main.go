package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"

	cdc "github.com/cmd-stream/codec-mus-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

const Addr = "127.0.0.1:9004"

func init() {
	assert.On = true
}

func main() {
	// Start server.
	var (
		greeter     = receiver.NewGreeter("Hello", "incredible", " ")
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		wgS         = &sync.WaitGroup{}
		server, err = utils.StartServer(Addr, serverCodec, greeter, wgS)
	)
	assert.EqualError(err, nil)

	// Create sender.
	clientCodec := cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	sender, err := MakeReconnectSender(clientCodec, connFactory{})
	assert.EqualError(err, nil)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)

	// Start the server again after some time.
	time.Sleep(time.Second)
	server, err = utils.StartServer(Addr, serverCodec, greeter, wgS)
	assert.EqualError(err, nil)

	// Wait for the sender clients to reconnect.
	time.Sleep(200 * time.Millisecond)

	wgR := &sync.WaitGroup{}
	// Send SayHelloCmd Сommand.
	wgR.Add(1)
	go func() {
		defer wgR.Done()
		var (
			sayHelloCmd  = cmds.NewSayHelloCmd("world")
			wantGreeting = results.Greeting("Hello world")
		)
		err = utils.Exchange(sayHelloCmd, wantGreeting, sender)
		assert.EqualError(err, nil)
	}()
	wgR.Wait()

	// Close sender.
	err = utils.CloseSender(sender)
	assert.EqualError(err, nil)

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)
}

type connFactory struct{}

func (f connFactory) New() (conn net.Conn, err error) {
	time.Sleep(100 * time.Millisecond)
	conn, err = net.Dial("tcp", Addr)
	if err != nil {
		log.Println("ConnFactory: failed to create a connection")
	} else {
		log.Println("ConnFactory: success")
	}
	return
}
