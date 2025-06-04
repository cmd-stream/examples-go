package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/examples-go/hello-world/utils"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	srv "github.com/cmd-stream/cmd-stream-go/server"
	ccln "github.com/cmd-stream/core-go/client"
	csrv "github.com/cmd-stream/core-go/server"
	sndr "github.com/cmd-stream/sender-go"
)

type listenerAdapter struct {
	net.Listener
	l *net.TCPListener
}

func (a listenerAdapter) SetDeadline(tm time.Time) error {
	return a.l.SetDeadline(tm)
}

func StartServer[T any](addr string, codec srv.Codec[T], receiver T,
	wg *sync.WaitGroup) (server *csrv.Server, err error) {
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		return
	}
	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	adapter := listenerAdapter{tls.NewListener(l, &tlsConf), l.(*net.TCPListener)}
	return utils.StartServerWith(addr, codec, srv.NewInvoker(receiver), adapter, wg)
}

func MakeSender[T any](addr string, clientsCount int, codec cln.Codec[T]) (
	sender sndr.Sender[T], err error) {
	var connFactory cln.ConnFactoryFn = func() (net.Conn, error) {
		cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
		if err != nil {
			log.Fatalf("Server: loadkeys: %s", err)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert},
			InsecureSkipVerify: true}
		return tls.Dial("tcp", addr, &config)
	}
	group, err := cmdstream.MakeClientGroup(clientsCount, codec, connFactory,
		grp.WithClientOps[T](
			cln.WithCore(
				ccln.WithUnexpectedResultCallback(utils.ClientCallback),
			),
		),
	)
	if err != nil {
		return
	}
	sender = sndr.New(group)
	return
}
