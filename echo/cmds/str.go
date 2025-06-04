package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/echo/results"
)

type StrCmd string

func (c StrCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{}, proxy core.Proxy) (err error) {
	_, err = proxy.Send(seq, results.Echo(c))
	return
}
