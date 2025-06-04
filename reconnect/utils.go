package main

import (
	cmdstream "github.com/cmd-stream/cmd-stream-go"
	cln "github.com/cmd-stream/cmd-stream-go/client"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	ccln "github.com/cmd-stream/core-go/client"
	"github.com/cmd-stream/examples-go/hello-world/utils"
	sndr "github.com/cmd-stream/sender-go"
)

func MakeReconnectSender[T any](codec cln.Codec[T], connFactory cln.ConnFactory) (
	sender sndr.Sender[T], err error) {
	group, err := cmdstream.MakeClientGroup(1, codec, connFactory,
		grp.WithReconnect[T](),
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
