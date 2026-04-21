package main

import (
	"context"
	"fmt"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cdcproto "github.com/cmd-stream/codec-protobuf-go"
	"github.com/cmd-stream/examples-go/calc_protobuf/cmds"
	rcvr "github.com/cmd-stream/examples-go/calc_protobuf/receiver"
	"github.com/cmd-stream/examples-go/calc_protobuf/results"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		registry = cdcproto.NewRegistry(
			cdcproto.WithCmd[rcvr.Calc, *cmds.AddCmd](),
			cdcproto.WithCmd[rcvr.Calc, *cmds.SubCmd](),
			cdcproto.WithResult[rcvr.Calc, *results.CalcResult](),
		)
		serverCodec = cdcproto.NewServerCodecWith(registry)
		clientCodec = cdcproto.NewClientCodecWith(registry)

		// Alternatively, you can create codecs manually using reflection types:
		//
		// cmdTypes = []reflect.Type{
		// 	 reflect.TypeFor[*cmds.AddCmd](),
		// 	 reflect.TypeFor[*cmds.SubCmd](),
		// }
		// resultTypes = []reflect.Type{
		// 	 reflect.TypeFor[*results.CalcResult](),
		// }
		// serverCodec = cdcproto.NewServerCodec[rcvr.Calc](cmdTypes, resultTypes)
		// clientCodec = cdcproto.NewClientCodec[rcvr.Calc](cmdTypes, resultTypes)
	)

	// Start server.
	fmt.Printf("Starting server on %s...\n", addr)
	go func() {
		server, _ := cmdstream.NewServer(rcvr.NewCalc(), serverCodec)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make sender.
	fmt.Println("Initializing sender and connecting...")
	sender, _ := cmdstream.NewSender(addr, clientCodec)

	// Send AddCmd.
	res, _ := sender.Send(context.Background(), &cmds.AddCmd{A: 2, B: 3})
	fmt.Printf("Sending AddCmd(2, 3)... Result: %v\n", res.(*results.CalcResult).R)

	// Send SubCmd.
	res, _ = sender.Send(context.Background(), &cmds.SubCmd{A: 8, B: 4})
	fmt.Printf("Sending SubCmd(8, 4)... Result: %v\n", res.(*results.CalcResult).R)
}
