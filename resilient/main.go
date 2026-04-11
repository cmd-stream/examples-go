package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	cdc "github.com/cmd-stream/codec-json-go"
	assert "github.com/ymz-ncnk/assert/panic"
	"github.com/ymz-ncnk/circbrk-go"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		cmdTypes = []reflect.Type{
			reflect.TypeFor[Message](),
		}
		resultTypes = []reflect.Type{
			reflect.TypeFor[Message](),
		}
		serverCodec = cdc.NewServerCodec[struct{}](cmdTypes, resultTypes)
		clientCodec = cdc.NewClientCodec[struct{}](cmdTypes, resultTypes)
	)

	// Start server.
	server, err := cmdstream.NewServer(struct{}{}, serverCodec)
	assert.EqualError(err, nil)
	defer server.Close()

	fmt.Printf("Starting server on %s...\n", addr)
	go func() {
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make resilient sender.
	fmt.Println("Initializing sender and connecting...")
	sender, err := MakeResilientSender(addr, clientCodec)
	assert.EqualError(err, nil)
	defer sender.Close()

	// Send message.
	msg1, err := sender.Send(context.Background(), Message("message1"))
	assert.EqualError(err, nil)
	fmt.Printf("Sending message1... Result: %v\n", msg1)

	// Close server.
	fmt.Printf("-- Closing server... --\n")
	server.Close()

	// Send messages to trigger (open) the circuit breaker.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = sender.Send(ctx2, Message("message2"))
	cancel2()
	fmt.Printf("Sending message2... Error: %v\n", err)

	ctx3, cancel3 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = sender.Send(ctx3, Message("message3"))
	cancel3()
	fmt.Printf("Sending message3... Error: %v\n", err)

	ctx4, cancel4 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = sender.Send(ctx4, Message("message4"))
	cancel4()
	fmt.Printf("Sending message4... Error: %v\n", err)

	fmt.Printf("Starting server on %s...\n", addr)
	go func() {
		server.ListenAndServe(addr)
	}()
	time.Sleep(2 * time.Second)

	// Send message.
	msg5, err := sender.Send(context.Background(), Message("message5"))
	assert.EqualError(err, nil)
	fmt.Printf("Sending message5... Result: %v\n", msg5)
}

func MakeResilientSender(addr string, clientCodec cln.Codec[struct{}]) (
	sndr.Sender[struct{}], error,
) {
	return cmdstream.NewSender(addr, clientCodec,
		sndr.WithClientsCount[struct{}](2),
		// Configure sender with circuit breaker hooks.
		sndr.WithSender(
			sndr.WithHooksFactory(
				hks.NewCircuitBreakerHooksFactory(
					circbrk.New(
						circbrk.WithWindowSize(4),
						circbrk.WithFailureRate(0.5),
						circbrk.WithOpenDuration(1*time.Second),
						circbrk.WithSuccessThreshold(1),
						circbrk.WithChangeStateCallback(
							func(state circbrk.State) {
								fmt.Printf("Circuit breaker state changed to: %v\n", state)
							},
						),
					),
					hks.NoopHooksFactory[struct{}]{},
				),
			),
		),
		// Configure group with reconnect and keepalive.
		sndr.WithGroup(
			grp.WithReconnect[struct{}](),
			grp.WithClient[struct{}](
				cln.WithKeepalive(
					dcln.WithKeepaliveIntvl(time.Second),
					dcln.WithKeepaliveTime(time.Second),
				),
			),
		),
	)
}
