package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/core-go"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/otel/cmds"
	otelcmd "github.com/cmd-stream/otelcmd-stream-go"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	ctx := context.Background()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	// Start server.
	var (
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		// Serializers for core.Cmd and core.Result interfaces allow building
		// a server codec.
		codec = cdc.NewServerCodec(results.ResultMUS, cmds.CmdMUS)
		wgS   = &sync.WaitGroup{}
	)
	tracerProvider, err := newServerTraceProvider()
	assert.EqualError(err, nil)
	defer func() {
		tracerProvider.Shutdown(ctx)
	}()

	server, err := StartServer(addr, codec, greeter, tracerProvider, wgS)
	assert.EqualError(err, nil)

	wgC := &sync.WaitGroup{}
	wgC.Add(1)
	go func() {
		SendCmds(addr, 90*time.Second)
		wgC.Done()
	}()

	<-time.NewTimer(30 * time.Second).C

	// Close server.
	log.Println("-- close server --")
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)

	<-time.NewTimer(30 * time.Second).C

	// Start server.
	log.Println("-- start server --")
	server, err = StartServer(addr, codec, greeter, tracerProvider, wgS)
	assert.EqualError(err, nil)

	wgC.Wait()

	// Close server.
	err = utils.CloseServer(server, wgS)
	assert.EqualError(err, nil)

}

func SendCmds(addr string, timeToWork time.Duration) {
	ctx := context.Background()

	// Create sender.
	tracerProvider, err := newClientTraceProvider()
	if err != nil {
		panic(err)
	}
	defer func() {
		tracerProvider.Shutdown(ctx)
	}()

	codec := cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	sender, err := MakeSender(addr, codec, tracerProvider)
	assert.EqualError(err, nil)

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Send Commands.
	timer1 := time.NewTimer(timeToWork)
	go func() {
		time.Sleep(30 * time.Millisecond)
		for {
			select {
			case <-timer1.C:
				goto Done
			default:
				cmd, wantGreeting := randCmd()
				_ = Exchange(ctx, cmd, sender, wantGreeting)
				// The server will be closed and started again, while we are trying
				// to send commands, so ignore errors here.
			}
		}
	Done:
		wg.Done()
	}()
	timer2 := time.NewTimer(timeToWork)
	go func() {
		time.Sleep(30 * time.Millisecond)
		for {
			select {
			case <-timer2.C:
				goto Done
			default:
				cmd, wantGreeting := randCmd()
				_ = Exchange(ctx, cmd, sender, wantGreeting)
				// The server will be closed and started again, while we are trying
				// to send commands, so ignore errors here.
			}
		}
	Done:
		wg.Done()
	}()

	wg.Wait()
	// Close sender.
	err = CloseSender(sender)
	assert.EqualError(err, nil)
}

func randCmd() (cmd core.Cmd[receiver.Greeter], wantGreeting results.Greeting) {
	var (
		n   = rand.Intn(5)
		str = gofakeit.Name()
	)
	switch n {
	case 0, 1, 2:
		cmd = otelcmd.NewTraceCmd(hwcmds.NewSayHelloCmd(str))
		wantGreeting = results.Greeting("Hello " + str)
	case 3:
		cmd = hwcmds.NewSayHelloCmd(str)
		wantGreeting = results.Greeting("Hello " + str)
	case 4, 5:
		cmd = otelcmd.NewTraceCmd(hwcmds.NewSayFancyHelloCmd(str))
		wantGreeting = results.Greeting("Hello incredible " + str)
	}
	return
}
