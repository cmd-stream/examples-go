package main

import (
	"context"
	"fmt"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/examples-go/calc/cmds"
	rcvr "github.com/cmd-stream/examples-go/calc/receiver"
	"github.com/cmd-stream/examples-go/calc/results"

	cdc "github.com/cmd-stream/codec-mus-stream-go"

	srv "github.com/cmd-stream/cmd-stream-go/server"

	sndr "github.com/cmd-stream/sender-go"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		invoker     = srv.NewInvoker[rcvr.Calc](rcvr.NewCalc())
		serverCodec = cdc.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = cdc.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
	)

	// Start server.
	go func() {
		server := cmdstream.MakeServer(serverCodec, invoker)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make sender.
	sender, _ := sndr.Make(addr, clientCodec)

	// Send AddCmd.
	sum, _ := sender.Send(context.Background(), cmds.AddCmd{A: 2, B: 3})
	fmt.Println(sum) // Output: 5

	// Send SubCmd.
	diff, _ := sender.Send(context.Background(), cmds.SubCmd{A: 8, B: 4})
	fmt.Println(diff) // Output: 4
}
