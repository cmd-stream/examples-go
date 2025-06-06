package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/cmd-stream/handler-go"
	sndr "github.com/cmd-stream/sender-go"
)

const (
	// CmdReceiveDuration specifies the timeout after which the server will
	// terminate the connection if no Command is received.
	CmdReceiveDuration    = 6 * time.Second
	CmdSendDuration       = time.Second
	ResultReceiveDuration = 3 * time.Second
)

// StartServer creates a listener, configures the server, and starts it
// in a separate goroutine.
func StartServer[T any](addr string, codec srv.Codec[T], receiver T,
	wg *sync.WaitGroup) (server *csrv.Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	invoker := srv.NewInvoker(receiver)
	return StartServerWith(addr, codec, invoker, l.(*net.TCPListener), wg)
}

func StartServerWith[T any](addr string, codec srv.Codec[T],
	invoker handler.Invoker[T],
	listener core.Listener,
	wg *sync.WaitGroup,
) (server *csrv.Server, err error) {
	var callback csrv.LostConnCallback = func(addr net.Addr, err error) {
		log.Printf("Server: lost connection to %v, cause %v\n", addr, err)
	}
	server = cmdstream.MakeServer(codec, invoker,
		// ServerInfo is optional and helps the client verify compatibility with the
		// server. It can identify supported commands or other server-specific
		// details. As a byte slice, it can store any arbitrary data.
		// srv.WithServerInfo(info)

		// Use Transport configuration to set the buffers size. If absent default
		// values from the bufio package are used.
		// srv.WithTransport(
		//   tcom.WithWriterBufSize(wsize),
		//   tcom.WithReaderBufSize(rsize)
		// )

		srv.WithCore(
			// WorkersCount determines the number of Workers, i.e., the number of
			// simultaneous connections to the server.
			csrv.WithWorkersCount(8),

			// LostConnCallback is useful for debugging, it is called by the server
			// when the connection to the client is lost.
			csrv.WithLostConnCallback(callback),
		),

		srv.WithHandler(
			// In a production environment, always set CmdReceiveTimeout. It allows
			// the server to close inactive client connections.
			handler.WithCmdReceiveDuration(CmdReceiveDuration),
			handler.WithAt(),
		),
	)

	wg.Add(1)
	go func(server *csrv.Server, l core.Listener, wg *sync.WaitGroup) {
		defer wg.Done()
		server.Serve(l)
	}(server, listener, wg)
	return
}

// CloseServer closes the server and waits for it to stop.
func CloseServer(server *csrv.Server, wg *sync.WaitGroup) (err error) {
	err = server.Close()
	if err != nil {
		return
	}

	// There is also the Server.Shutdown() method, which allows the server
	// to complete processing of already established connections without
	// accepting new ones.

	wg.Wait()
	return
}

var ClientCallback ccln.UnexpectedResultCallback = func(seq core.Seq, result core.Result) {
	log.Printf("Client: unexpected result: seq %v, result %v\n", seq, result)
}

func MakeSender[T any](addr string, clientsCount int, codec cln.Codec[T]) (
	sender sndr.Sender[T], err error) {
	var connFactory cln.ConnFactoryFn = func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
	group, err := cmdstream.MakeClientGroup(clientsCount, codec, connFactory,
		grp.WithClientOps[T](
			cln.WithCore(
				// UnexpectedResultCallback handles unexpected results from the
				// server. If you call Client.Forget(seq) for a Command, its results
				// will be handled by this function.
				ccln.WithUnexpectedResultCallback(ClientCallback),
				// Use Transport configuration to set the buffers size. If absent
				// default values from the bufio package are used.
				// ccln.WithTransport(
				//   tcom.WithWriterBufSize(...),
				//   tcom.WithReaderBufSize(...),
				// ),
			),
		),
	)
	if err != nil {
		// An error in this case indicates that the number of clients in the group is
		// less than the clientsCount parameter. If you're not satisfied with that,
		// don't forget to close the group:
		group.Close()
		return
	}
	sender = sndr.New(group)
	return
}

// CloseSender closes the sender and waits for it to stop.
func CloseSender[T any](sender sndr.Sender[T]) (err error) {
	// Sender uses several clients, and each client will be closed if:
	// - Sender.Close() is called.
	// - The server terminates the connection.

	err = sender.Close()
	if err != nil {
		return
	}
	// Clients receive Results from the server in the background, so we
	// have to wait until they stop.
	select {
	case <-time.NewTimer(time.Second).C:
		return errors.New("can't close the sender, cause timeout exceeded")
	case <-sender.Done():
		return
	}
}

// Exchange sends a Command and checks whether the received greeting
// matches the expected value.
func Exchange[T any, R interface{ String() string }](cmd core.Cmd[T],
	wantGreeting R,
	sender sndr.Sender[T],
) (err error) {
	// Send the Command.
	var (
		ctx, cancel = context.WithTimeout(context.Background(), ResultReceiveDuration)
		deadline    = time.Now().Add(CmdSendDuration)
	)
	defer cancel()

	result, err := sender.SendWithDeadline(ctx, cmd, deadline)
	if err != nil {
		return
	}
	greeting := result.(R)
	if greeting.String() != wantGreeting.String() {
		return fmt.Errorf("unexpected greeting, want %v actual %v", wantGreeting,
			greeting)
	}
	return
}
