package main

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
)

// Message implements both core.Cmd and core.Result interfaces.
type Message string

func (s Message) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{}, proxy core.Proxy,
) (err error) {
	_, err = proxy.Send(seq, s)
	return
}

func (s Message) LastOne() bool {
	return true
}
