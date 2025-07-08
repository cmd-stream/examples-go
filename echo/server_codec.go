package main

import (
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/echo/cmds"
	"github.com/cmd-stream/examples-go/echo/results"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// ServerCodec used to initialize the server.
type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w transport.Writer) (n int, err error) {
	n, err = ord.String.Marshal(string(result.(results.EchoResult)), w)
	return
}

func (c ServerCodec) Decode(r transport.Reader) (EchoCmd core.Cmd[struct{}], n int,
	err error) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	EchoCmd = cmds.EchoCmd(str)
	return
}
