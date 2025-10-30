package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	rcvr "github.com/cmd-stream/examples-go/calc_protobuf/receiver"
	"github.com/cmd-stream/examples-go/calc_protobuf/results"
)

func (c *SubCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	calc rcvr.Calc, proxy core.Proxy,
) (err error) {
	result := results.CalcResult{R: calc.Sub(c.A, c.B)}
	_, err = proxy.Send(seq, &result)
	return
}
