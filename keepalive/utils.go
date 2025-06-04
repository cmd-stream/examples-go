package main

import (
	"net"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	ccln "github.com/cmd-stream/core-go/client"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	sndr "github.com/cmd-stream/sender-go"
)

const (
	KeepaliveTime  = 200 * time.Millisecond
	KeepaliveIntvl = 200 * time.Millisecond
)

func MakeKeepaliveSender[T any](addr string, clientsCount int,
	codec cln.Codec[T]) (sender sndr.Sender[T], err error) {
	var connFactory cln.ConnFactoryFn = func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
	group, err := cmdstream.MakeClientGroup(clientsCount, codec, connFactory,
		grp.WithClientOps[T](
			cln.WithCore(
				ccln.WithUnexpectedResultCallback(utils.ClientCallback),
			),
			cln.WithKeepalive(
				dcln.WithKeepaliveTime(KeepaliveTime),
				dcln.WithKeepaliveIntvl(KeepaliveIntvl),
			),
		),
	)
	if err != nil {
		return
	}
	sender = sndr.New(group)
	return
}
