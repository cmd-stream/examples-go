package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
)

// SayFancyHelloCmd implements core.Cmd and exts.MarshallerTypedMUS interfaces.
// Produces greetings like "Hello incredible world".
type SayFancyHelloCmd struct {
	Str string `json:"str"`
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	greeter receiver.Greeter, proxy core.Proxy,
) (err error) {
	// SayFancyHelloCmd differs from SayHelloCmd in the way it uses the
	// Receiver.
	var (
		greeting = results.Greeting(
			greeter.Join(greeter.Interjection(), greeter.Adjective(), c.Str),
		)
		deadline = at.Add(CmdExecDuration)
	)
	_, err = proxy.SendWithDeadline(seq, greeting, deadline)
	return
}
