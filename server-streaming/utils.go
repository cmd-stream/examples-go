package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/server-streaming/results"
	sndr "github.com/cmd-stream/sender-go"
)

func Exchange[T any](cmd core.Cmd[T], wantGreetings []results.Greeting,
	sender sndr.Sender[T]) (err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(),
			utils.ResultReceiveDuration)
		deadline = time.Now().Add(utils.CmdSendDuration)
	)
	defer cancel()
	var handler sndr.ResultHandlerFn = func(result core.Result, err error) error {
		greeting := result.(results.Greeting)
		if !reflect.DeepEqual(greeting, wantGreetings[0]) {
			return fmt.Errorf("unexpected greeting, want %v actual %v",
				wantGreetings[0], greeting)
		}
		return nil
	}
	return sender.SendMultiWithDeadline(ctx, cmd, len(wantGreetings), handler,
		deadline)
}
