package utils

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/results"
	sndr "github.com/cmd-stream/sender-go"
)

const (
	// CmdReceiveDuration specifies the timeout after which the server will
	// terminate the connection if no Command is received.
	CmdReceiveDuration    = 6 * time.Second
	CmdSendDuration       = time.Second
	ResultReceiveDuration = 3 * time.Second
)

func SendCmd[T any](cmd core.Cmd[T], sender sndr.Sender[T]) (
	greeting results.Greeting, err error,
) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(),
			ResultReceiveDuration)
		deadline = time.Now().Add(CmdSendDuration)
	)
	defer cancel()
	result, err := sender.SendWithDeadline(ctx, cmd, deadline)
	if err != nil {
		return
	}
	greeting = result.(results.Greeting)
	return
}
