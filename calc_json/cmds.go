package main

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// AddCmd defines the client request for an addition operation. As a Command,
// it implements the core.Command interface via its Exec method.
type AddCmd struct {
	A, B int
}

func (c AddCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	calc Calc, proxy core.Proxy,
) (err error) {
	result := CalcResult(calc.Add(c.A, c.B))
	_, err = proxy.Send(seq, result)
	return
}

// SubCmd defines the client request for a substraction operation. As a
// Command, it implements the core.Command interface via its Exec method.
type SubCmd struct {
	A, B int
}

func (c SubCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	calc Calc, proxy core.Proxy,
) (err error) {
	result := CalcResult(calc.Sub(c.A, c.B))
	_, err = proxy.Send(seq, result)
	return
}
