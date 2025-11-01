package main

import (
	"log"
	"net"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	codecmus "github.com/cmd-stream/codec-mus-stream-go"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/handler-go"

	srv "github.com/cmd-stream/cmd-stream-go/server"

	ccln "github.com/cmd-stream/core-go/client"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	sndr "github.com/cmd-stream/sender-go"
	assert "github.com/ymz-ncnk/assert/panic"

	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
)

func init() {
	assert.On = true
}

func main() {
	const addr = "127.0.0.1:9000"
	var (
		greeter = rcvr.NewGreeter("Hello", "incredible", " ")
		invoker = srv.NewInvoker[rcvr.Greeter](greeter)
		// Serializers of core.Cmd and core.Result interfaces allow building
		// server/client codecs.
		serverCodec = codecmus.NewServerCodec(cmds.CmdMUS, results.ResultMUS)
		clientCodec = codecmus.NewClientCodec(cmds.CmdMUS, results.ResultMUS)
		// // Alternative JSON codecs, require github.com/cmd-stream/codec-json-go:
		// cmdTypes = []reflect.Type{
		//   reflect.TypeFor[cmds.SayHelloCmd](),
		//   reflect.TypeFor[cmds.SayFancyHelloCmd](),
		// }
		// resultTypes = []reflect.Type{
		//   reflect.TypeFor[results.Greeting](),
		// }
		// serverCodec = codecjson.NewServerCodec[rcvr.Greeter](cmdTypes, resultTypes)
		// clientCodec = codecjson.NewClientCodec[rcvr.Greeter](cmdTypes, resultTypes)
		wgS = &sync.WaitGroup{}
	)

	server := MakeServer(serverCodec, invoker)
	wgS.Add(1)
	go func() {
		server.ListenAndServe(addr)
		wgS.Done()
	}()
	time.Sleep(100 * time.Millisecond)

	// Instead of an asynchronious client we will use Sender, that is build on
	// the group of clients.
	sender, err := MakeSender(addr, clientCodec)
	assert.EqualError(err, nil)
	SendCmds(sender)

	// Sender is built on one or more clients that receive results from the
	// server in the background, so we need to wait for them to finish.
	err = sender.CloseAndWait(time.Second)
	assert.EqualError(err, nil)
	// In addition to Server.Close(), there is also the Server.Shutdown() method,
	// which allows the server to complete processing of already established
	// connections without accepting new ones.
	err = server.Close()
	assert.EqualError(err, nil)
	wgS.Wait()
}

func MakeServer(codec srv.Codec[rcvr.Greeter],
	invoker handler.Invoker[rcvr.Greeter],
) *csrv.Server {
	return cmdstream.MakeServer(codec, invoker,
		// // ServerInfo is optional and helps the client verify compatibility with the
		// // server. It can identify supported commands or other server-specific
		// // etails. As a byte slice, it can store any arbitrary data.
		// srv.WithServerInfo(...),

		srv.WithCore(
			// WorkersCount determines the number of Workers, i.e., the number of
			// simultaneous connections to the server.
			csrv.WithWorkersCount(2),

			// LostConnCallback is useful for debugging, it is called by the server
			// when the connection to the client is lost.
			csrv.WithLostConnCallback(func(addr net.Addr, err error) {
				log.Printf("server: lost connection to %v, cause %v\n", addr, err)
			}),
		),
		srv.WithHandler(
			// In a production environment, always set CmdReceiveTimeout. It allows
			// the server to close inactive client connections.
			handler.WithCmdReceiveDuration(utils.CmdReceiveDuration),
			// The Commands being sent uses the 'at' parameter, so enable it.
			handler.WithAt(),
		),

		// // Use Transport configuration to set the buffers size. If absent default
		// // values from the bufio package are used.
		// srv.WithTransport(
		//  transport.WithWriterBufSize(...),
		//  transport.WithReaderBufSize(...),
		// ),
	)
}

func MakeSender(addr string, codec cln.Codec[rcvr.Greeter]) (
	sender sndr.Sender[rcvr.Greeter], err error,
) {
	return sndr.Make(addr, codec,
		sndr.WithGroup(
			grp.WithClient[rcvr.Greeter](
				// // Optional ServerInfo.
				// cln.WithServerInfo(...),

				// UnexpectedResultCallback handles unexpected results from the
				// server. If you call Client.Forget(seq) for a Command, its results
				// will be handled by this function. The default callback logs a message to
				// the standard logger.
				cln.WithCore(
					ccln.WithUnexpectedResultCallback(func(seq core.Seq, result core.Result) {
						log.Printf("client: unexpected result: seq %v, result %v\n", seq, result)
					}),
				),

				// // Use Transport configuration to set the buffers size. If absent default
				// // values from the bufio package are used.
				// cln.WithTransport(
				//  transport.WithWriterBufSize(...),
				//  transport.WithReaderBufSize(...),
				// ),
			),
		),
		// // By using an optional hooks factory, the sender creates and applies
		// // a new set of hooks each time it sends a Command. Hooks let you
		// // extend or modify the sender's behavior, for example, by adding
		// // logging, metrics, retries, or circuit breaker logic.
		// sndr.WithSender(
		//  sndr.WithHooksFactory[rcvr.Greeter](...),
		// ),
		sndr.WithClientsCount[rcvr.Greeter](2),
	)

	// // If you want a sender that takes care of connection management for you,
	// // use this setup. It configures the sender to automatically handle
	// // keepalives, reconnects, and circuit breaker behavior, so as long as
	// // the server is alive, you’ll stay connected.
	// var (
	// 	cb = circbrk.New(circbrk.WithWindowSize(20),
	// 		circbrk.WithFailureRate(0.5),
	// 		circbrk.WithOpenDuration(30*time.Second),
	// 		circbrk.WithSuccessThreshold(2),
	// 	)
	// 	hooksFactory = hks.NewCircuitBreakerHooksFactory(cb,
	// 		hks.NoopHooksFactory[rcvr.Greeter]{})
	// )
	// return sndr.Make(addr, codec,
	// 	sndr.WithGroup(
	// 		grp.WithReconnect[rcvr.Greeter](),
	// 		grp.WithClient[rcvr.Greeter](
	// 			cln.WithCore(
	// 				ccln.WithUnexpectedResultCallback(
	// 					func(seq core.Seq, result core.Result) {
	// 						log.Printf("client: unexpected result: seq %v, result %v\n", seq, result)
	// 					},
	// 				),
	// 			),
	// 			cln.WithKeepalive(
	// 				dcln.WithKeepaliveTime(30*time.Second),
	// 				dcln.WithKeepaliveIntvl(10*time.Second),
	// 			),
	// 		),
	// 	),
	// 	sndr.WithSender(
	// 		sndr.WithHooksFactory(hooksFactory),
	// 	),
	// 	sndr.WithClientsCount[rcvr.Greeter](2),
	// )
}

func SendCmds(sender sndr.Sender[rcvr.Greeter]) {
	wg := sync.WaitGroup{}

	// Send SayHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.SayHelloCmd{Str: "world"}
			want = results.Greeting("Hello world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting, want)
		wg.Done()
	}()

	// Send SayFancyHelloCmd.
	wg.Add(1)
	go func() {
		var (
			cmd  = cmds.SayFancyHelloCmd{Str: "world"}
			want = results.Greeting("Hello incredible world")
		)
		greeting, err := utils.SendCmd(cmd, sender)
		assert.EqualError(err, nil)
		assert.Equal(greeting, want)
		wg.Done()
	}()

	wg.Wait()
}
