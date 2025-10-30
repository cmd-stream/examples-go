package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	rcvr "github.com/cmd-stream/examples-go/calc_json/receiver"
	"github.com/cmd-stream/examples-go/calc_json/results"
)

// SubCmd defines the client request for a substraction operation. As a
// Command, it implements the core.Command interface via its Exec method.
type SubCmd struct {
	A, B int
}

func (c SubCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	calc rcvr.Calc, proxy core.Proxy,
) (err error) {
	result := results.CalcResult(calc.Sub(c.A, c.B))
	_, err = proxy.Send(seq, result)
	return
}
