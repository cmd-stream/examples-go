package main

import (
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/echo/cmds"
	"github.com/cmd-stream/examples-go/echo/results"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// ClientCodec used to initialize the client.
type ClientCodec struct{}

func (c ClientCodec) Encode(cmd core.Cmd[struct{}], w transport.Writer) (
	n int, err error) {
	_, err = ord.String.Marshal(string(cmd.(cmds.StrCmd)), w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result core.Result, n int,
	err error) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	result = results.Echo(str)
	return
}
