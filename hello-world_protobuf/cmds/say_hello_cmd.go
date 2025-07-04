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

// NewSayHelloCmd creates a new SayHelloCmd.
func NewSayHelloCmd(str string) SayHelloCmd {
	return SayHelloCmd{
		SayHelloData: &SayHelloData{Str: str},
	}
}

// SayHelloCmd implements core.Cmd and Marshaller interfaces.
type SayHelloCmd struct {
	*SayHelloData
}

func (c SayHelloCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	greeter receiver.Greeter,
	proxy core.Proxy,
) (err error) {
	var (
		str      = greeter.Join(greeter.Interjection(), c.Str)
		result   = results.NewGreeting(str)
		deadline = at.Add(utils.CmdSendDuration)
	)
	_, err = proxy.SendWithDeadline(seq, result, deadline)
	return
}

func (c SayHelloCmd) MarshalTypedProtobuf(w muss.Writer) (n int,
	err error) {
	return SayHelloCmdDTS.Marshal(c, w)
}

func (c SayHelloCmd) SizeTypedProtobuf() (size int) {
	return SayHelloCmdDTS.Size(c)
}
