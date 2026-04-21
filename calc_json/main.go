package main

import (
	"context"
	"fmt"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cdcjson "github.com/cmd-stream/codec-json-go"
)

func main() {
	const addr = "127.0.0.1:9000"
	var (
		registry = cdcjson.NewRegistry(
			cdcjson.WithCmd[Calc, AddCmd](),
			cdcjson.WithCmd[Calc, SubCmd](),
			cdcjson.WithResult[Calc, CalcResult](),
		)
		serverCodec = cdcjson.NewServerCodecWith(registry)
		clientCodec = cdcjson.NewClientCodecWith(registry)

		// Alternatively, you can create codecs manually using reflection types:
		//
		// cmdTypes = []reflect.Type{
		// 	 reflect.TypeFor[AddCmd](),
		// 	 reflect.TypeFor[SubCmd](),
		// }
		// resultTypes = []reflect.Type{
		// 	 reflect.TypeFor[CalcResult](),
		// }
		// serverCodec = cdcjson.NewServerCodec[Calc](cmdTypes, resultTypes)
		// clientCodec = cdcjson.NewClientCodec[Calc](cmdTypes, resultTypes)
	)

	// Start server.
	fmt.Printf("Starting server on %s...\n", addr)
	go func() {
		server, _ := cmdstream.NewServer(NewCalc(), serverCodec)
		server.ListenAndServe(addr)
	}()
	time.Sleep(100 * time.Millisecond)

	// Make sender.
	fmt.Println("Initializing sender and connecting...")
	sender, _ := cmdstream.NewSender(addr, clientCodec)

	// Send AddCmd.
	sum, _ := sender.Send(context.Background(), AddCmd{A: 2, B: 3})
	fmt.Printf("Sending AddCmd(2, 3)... Result: %v\n", sum)

	// Send SubCmd.
	diff, _ := sender.Send(context.Background(), SubCmd{A: 8, B: 4})
	fmt.Printf("Sending SubCmd(8, 4)... Result: %v\n", diff)
}
