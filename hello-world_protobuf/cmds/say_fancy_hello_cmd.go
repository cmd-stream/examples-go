package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	receiver "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	"github.com/cmd-stream/examples-go/hello-world_protobuf/results"
	muss "github.com/mus-format/mus-stream-go"
)

// NewSayFancyHelloCmd creates a new SayFancyHelloCmd.
func NewSayFancyHelloCmd(str string) SayFancyHelloCmd {
	return SayFancyHelloCmd{
		SayFancyHelloData: &SayFancyHelloData{Str: str},
	}
}

// SayFancyHelloCmd implements core.Cmd and Marshaller interfaces.
type SayFancyHelloCmd struct {
	*SayFancyHelloData
}

func (c SayFancyHelloCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	greeter receiver.Greeter,
	proxy core.Proxy,
) (err error) {
	var (
		str      = greeter.Join(greeter.Interjection(), greeter.Adjective(), c.Str)
		result   = results.NewGreeting(str)
		deadline = at.Add(utils.CmdSendDuration)
	)
	_, err = proxy.SendWithDeadline(seq, result, deadline)
	return
}

func (c SayFancyHelloCmd) MarshalTypedProtobuf(w muss.Writer) (n int,
	err error) {
	return SayFancyHelloCmdDTS.Marshal(c, w)
}

func (c SayFancyHelloCmd) SizeTypedProtobuf() (size int) {
	return SayFancyHelloCmdDTS.Size(c)
}
