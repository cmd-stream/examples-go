package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	muss "github.com/mus-format/mus-stream-go"
)

// NewSayFancyHelloCmd creates a new SayFancyHelloCmd.
func NewSayFancyHelloCmd(str string) SayFancyHelloCmd {
	return SayFancyHelloCmd{str}
}

// SayFancyHelloCmd implements the core.Cmd[Greeter] interface and produces
// greetings like "Hello incredible world".
type SayFancyHelloCmd struct {
	str string
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	greeter receiver.Greeter, proxy core.Proxy) (err error) {
	// SayFancyHelloCmd differs from SayHelloCmd in the way it uses the
	// Receiver.
	var (
		greeting = results.Greeting(
			greeter.Join(greeter.Interjection(), greeter.Adjective(), c.str),
		)
		deadline = at.Add(CmdExecDuration)
	)
	_, err = proxy.SendWithDeadline(seq, greeting, deadline)
	return
}

func (c SayFancyHelloCmd) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return SayFancyHelloCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloCmd) SizeTypedMUS() (size int) {
	return SayFancyHelloCmdDTS.Size(c)
}
