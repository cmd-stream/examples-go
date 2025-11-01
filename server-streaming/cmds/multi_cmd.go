package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/server-streaming/results"
	"github.com/mus-format/mus-stream-go"
)

// NewSayFancyHelloMultiCmd creates a new SayFancyHelloMultiCmd.
func NewSayFancyHelloMultiCmd(str string) SayFancyHelloMultiCmd {
	return SayFancyHelloMultiCmd{str}
}

// SayFancyHelloMultiCmd implements core.Cmd and ext.MarshallerTypedMUS
// interfaces.
//
// We have to define MarshalTypedMUS and SizeTypedMUS methods (implement the
// MarshallerTypedMUS interface) because the core.Cmd interface
// serialization code was generated with introps.WithRegisterMarshaller().
type SayFancyHelloMultiCmd struct {
	str string
}

func (c SayFancyHelloMultiCmd) Exec(ctx context.Context, seq core.Seq,
	at time.Time,
	greeter receiver.Greeter,
	proxy core.Proxy,
) (err error) {
	var (
		deadline = at.Add(utils.CmdSendDuration)
		greeting = results.NewGreeting(greeter.Interjection(), false)
	)
	_, err = proxy.SendWithDeadline(seq, greeting, deadline)
	if err != nil {
		return
	}
	greeting = results.NewGreeting(greeter.Adjective(), false)
	_, err = proxy.Send(seq, greeting)
	if err != nil {
		return
	}
	greeting = results.NewGreeting(c.str, true)
	_, err = proxy.Send(seq, greeting)
	return
}

func (c SayFancyHelloMultiCmd) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return SayFancyHelloMultiCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloMultiCmd) SizeTypedMUS() (size int) {
	return SayFancyHelloMultiCmdDTS.Size(c)
}
