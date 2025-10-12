package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	rcvr "github.com/cmd-stream/examples-go/calc/receiver"
	"github.com/cmd-stream/examples-go/calc/results"
)

// AddCmd defines the client request for an addition operation. As a Command,
// it implements the core.Command interface via its Exec method.
type AddCmd struct {
	A, B int
}

func (c AddCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	calc rcvr.Calc, proxy core.Proxy,
) (err error) {
	result := results.CalcResult(calc.Add(c.A, c.B))
	_, err = proxy.Send(seq, result)
	return
}
