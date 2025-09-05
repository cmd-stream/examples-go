package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	muss "github.com/mus-format/mus-stream-go"
)

// NewSayHelloCmd creates a new SayHelloCmd.
func NewSayHelloCmd(str string) SayHelloCmd {
	return SayHelloCmd{str}
}

// SayHelloCmd implements core.Cmd and exts.MarshallerTypedMUS interfaces.
// Produces greetings like "Hello world".
type SayHelloCmd struct {
	str string
}

func (c SayHelloCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
	greeter receiver.Greeter, proxy core.Proxy,
) (err error) {
	var (
		greeting = results.Greeting(
			greeter.Join(greeter.Interjection(), c.str),
		)
		// Limiting the execution time of a Command on the server is
		// considered a good practice that can be achieved with a deadline.
		deadline = at.Add(CmdExecDuration)
	)

	// A Command can behave in various ways:
	// 1. It can send back only one Result:
	//
	//      return proxy.SendWithDeadline(seq, greeting, deadline)
	//
	//    Note: The deadline is applied to the entire connection. This means
	//    it will also affect subsequent commands unless they update it with
	//    their own value.
	//
	//    So if one Command uses the Proxy.SendWithDeadline() method, all
	//    others should do the same. Mixing Proxy.Send() and
	//    Proxy.SendWithDeadline() can result in unpredictable behavior
	//    due to unintended deadline propagation.
	//
	//    To cancel the deadline, use time.Time{}:
	//
	//      return proxy.SendWithDeadline(time.Time{}, seq, result)
	//
	// 2. It can perform context-related tasks:
	//
	//      ownCtx, cancel := context.WithDeadline(ctx, deadline)
	//      // Use ownCtx to perform a context-related task.
	//      ...
	//      return proxy.SendWithDeadline(deadline, seq, result)
	//
	// 3. It can send back multiple results (server streaming):
	//
	//      err = proxy.SendWithDeadline(deadline, seq, result1)
	//      if err != nil {
	//         return
	//      }
	//      return proxy.Send(seq, result2)
	//
	//    All results except the first one are sent back using the
	//    Proxy.Send() method.
	//
	// Regardless of the case, the final Result should have
	// Result.LastOne() == true.

	// As you can see, the current Command sends back only one Result.
	_, err = proxy.SendWithDeadline(seq, greeting, deadline)
	return
}

func (c SayHelloCmd) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return SayHelloCmdDTS.Marshal(c, w) // The Command will be marshalled as
	// 'DTM + command data'.
}

func (c SayHelloCmd) SizeTypedMUS() (size int) {
	return SayHelloCmdDTS.Size(c)
}
