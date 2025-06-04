package cmds

import (
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	otelcmd "github.com/cmd-stream/otelcmd-stream-go"
)

type TraceSayHelloCmd = otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd]
type TraceSayFancyHelloCmd = otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd]
