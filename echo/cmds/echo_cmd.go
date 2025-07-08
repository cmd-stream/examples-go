package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/echo/results"
)

type EchoCmd string

func (c EchoCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{}, proxy core.Proxy) (err error) {
	_, err = proxy.Send(seq, results.EchoResult(c))
	return
}
