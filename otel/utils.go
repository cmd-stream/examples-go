package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	otelcmd "github.com/cmd-stream/otelcmd-stream-go"
	sndr "github.com/cmd-stream/sender-go"
	hks "github.com/cmd-stream/sender-go/hooks"
	"github.com/ymz-ncnk/circbrk-go"
	"go.opentelemetry.io/otel/trace"
)

const (
	CircuitBreakerWindowSize       = 8
	CircuitBreakerFailureRate      = 0.5
	CircuitBreakerOpenDuration     = 6 * time.Second
	CircuitBreakerSuccessThreshold = 3
)

func StartServer[T any](addr string, codec srv.Codec[T], receiver T,
	tracerProvider trace.TracerProvider,
	wg *sync.WaitGroup,
) (server *csrv.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	invoker := otelcmd.NewInvoker(srv.NewInvoker(receiver),
		otelcmd.WithServerAddr[T](l.Addr()),
		otelcmd.WithTracerProvider[T](tracerProvider),
	)
	return utils.StartServerWith(addr, codec, invoker, l.(*net.TCPListener), wg)
}

func MakeSender[T any](addr string, codec cln.Codec[T],
	tracerProvider trace.TracerProvider) (sender sndr.Sender[T], err error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	var (
		connFactory cln.ConnFactoryFn = func() (net.Conn, error) {
			return net.Dial("tcp", addr)
		}
		hooksFactory = MakeSenderHooksFactory[T](serverAddr, tracerProvider)
	)
	group, err := cmdstream.MakeClientGroup(2, codec, connFactory,
		grp.WithReconnect[T](),
		grp.WithClientOps[T](
			cln.WithCore(
				ccln.WithUnexpectedResultCallback(utils.ClientCallback),
			),
		),
	)
	if err != nil {
		group.Close()
		return
	}
	sender = sndr.New(group, sndr.WithHooksFactory(hooksFactory))
	return
}

func MakeSenderHooksFactory[T any](serverAddr net.Addr,
	tracerProvider trace.TracerProvider) hks.CircuitBreakerHooksFactory[T] {
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

// CloseSender closes the sender and waits for it to stop.
func CloseSender[T any](sender sndr.Sender[T]) (err error) {
	// Generally, the client will be closed if:
	// - Client.Close() is called.
	// - The server terminates the connection.
	//
	// In both cases, all uncompleted Commands will receive
	// AsyncResult.Error() == Client.Err(), where Client.Err() returns a
	// connection error.

	err = sender.Close()
	if err != nil {
		return
	}
	// The sender receives Results from the server in the background, so we
	// have to wait until it stops.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("can't close the sender, cause timeout exceeded")
	case <-sender.Done():
		return
	}
}

// Exchange sends a Command and checks whether the received greeting
// matches the expected value.
func Exchange[T any, R interface{ String() string }](ctx context.Context,
	cmd core.Cmd[T],
	sender sndr.Sender[T],
	wantGreeting R,
) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result, err := sender.Send(ctx, cmd)
	if err != nil {
		return
	}
	greeting := result.(R)
	if greeting.String() != wantGreeting.String() {
		return fmt.Errorf("unexpected greeting, want %v actual %v", wantGreeting,
			greeting)
	}
	return
}
