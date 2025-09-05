package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	cmdstream "github.com/cmd-stream/cmd-stream-go"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	cdc "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/core-go"
	csrv "github.com/cmd-stream/core-go/server"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/otel/cmds"
	"github.com/cmd-stream/handler-go"
	otelcmd "github.com/cmd-stream/otelcmd-stream-go"
	sndr "github.com/cmd-stream/sender-go"
	hks "github.com/cmd-stream/sender-go/hooks"
	assert "github.com/ymz-ncnk/assert/panic"
	"github.com/ymz-ncnk/circbrk-go"

	// "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	CircuitBreakerWindowSize       = 8
	CircuitBreakerFailureRate      = 0.5
	CircuitBreakerOpenDuration     = 6 * time.Second
	CircuitBreakerSuccessThreshold = 3
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"

	// Set up OpenTelemetry -----------------------------------------------------
	ctx := context.Background()
	otelShutdown, err := setupOTelSDK(ctx)
	assert.EqualError(err, nil)
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err := otelShutdown(ctx)
		assert.EqualError(err, nil)
	}()
	serverTracerProvider, err := newServerTraceProvider()
	assert.EqualError(err, nil)
	defer func() {
		err := serverTracerProvider.Shutdown(ctx)
		assert.EqualError(err, nil)
	}()
	clientTracerProvider, err := newClientTraceProvider()
	assert.EqualError(err, nil)
	defer func() {
		err := clientTracerProvider.Shutdown(ctx)
		assert.EqualError(err, nil)
	}()
	// --------------------------------------------------------------------------

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	assert.EqualError(err, nil)

	var (
		greeter = receiver.NewGreeter("Hello", "incredible", " ")
		invoker = otelcmd.NewInvoker(srv.NewInvoker(greeter),
			otelcmd.WithServerAddr[receiver.Greeter](tcpAddr),
			otelcmd.WithTracerProvider[receiver.Greeter](serverTracerProvider),
		)
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		wgS         = &sync.WaitGroup{}
	)

	// Make server.
	server := cmdstream.MakeServer(serverCodec, invoker,
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
	hooksFactory := MakeSenderHooksFactory[receiver.Greeter](tcpAddr,
		clientTracerProvider)
	sender, err := sndr.Make(addr, clientCodec,
		sndr.WithGroup(
			grp.WithReconnect[receiver.Greeter](),
		),
		sndr.WithSender(
			sndr.WithHooksFactory(hooksFactory),
		),
		sndr.WithClientsCount[receiver.Greeter](2),
	)
	assert.EqualError(err, nil)
	// Send Commands.
	SendCmdsWithServerRestart(addr, sender, server, wgS)

	// Close sender.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// Close server.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func SendCmdsWithServerRestart(addr string,
	sender sndr.Sender[receiver.Greeter],
	server *csrv.Server,
	wgS *sync.WaitGroup,
) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		SendCmds(sender, 90*time.Second)
		wg.Done()
	}()

	<-time.NewTimer(30 * time.Second).C

	// Close server.
	log.Println("-- close server --")
	err := server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()

	<-time.NewTimer(30 * time.Second).C

	// Start server.
	log.Println("-- start server --")
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()

	wg.Wait()
}

func SendCmds(sender sndr.Sender[receiver.Greeter], timeToWork time.Duration) {
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
				cmd, want := RandCmd()
				greeting, err := utils.SendCmd(cmd, sender)
				// The server will be closed and started again, while we are trying
				// to send Commands, so ignore errors here.
				if err == nil {
					assert.Equal(greeting, want)
				}
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
				cmd, want := RandCmd()
				greeting, err := utils.SendCmd(cmd, sender)
				// The server will be closed and started again, while we are trying
				// to send Commands, so ignore errors here.
				if err == nil {
					assert.Equal(greeting, want)
				}
			}
		}
	Done:
		wg.Done()
	}()

	wg.Wait()
}

func RandCmd() (cmd core.Cmd[receiver.Greeter], wantGreeting results.Greeting) {
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

func MakeSenderHooksFactory[T any](serverAddr net.Addr,
	tracerProvider trace.TracerProvider,
) hks.CircuitBreakerHooksFactory[T] {
	var (
		cb = circbrk.New(circbrk.WithWindowSize(CircuitBreakerWindowSize),
			circbrk.WithFailureRate(CircuitBreakerFailureRate),
			circbrk.WithOpenDuration(CircuitBreakerOpenDuration),
			circbrk.WithSuccessThreshold(CircuitBreakerSuccessThreshold),
			circbrk.WithChangeStateCallback(
				func(state circbrk.State) {
					log.Printf("CircuitBreaker: %v", state)
				},
			),
		)
		otelHooksFactory = otelcmd.NewHooksFactory(
			otelcmd.WithServerAddr[T](serverAddr),
			otelcmd.WithTracerProvider[T](tracerProvider),
		)
	)
	return hks.NewCircuitBreakerHooksFactory(cb, otelHooksFactory)
}
