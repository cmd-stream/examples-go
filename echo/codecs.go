package main

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/transport"
	examples "github.com/cmd-stream/examples-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w transport.Writer) (
	n int, err error,
) {
	_, err = ord.String.Marshal(string(result.(examples.Message)), w)
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd core.Cmd[struct{}], n int,
	err error,
) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	cmd = examples.Message(str)
	return
}

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd core.Cmd[struct{}], w transport.Writer) (
	n int, err error,
) {
	_, err = ord.String.Marshal(string(cmd.(examples.Message)), w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result core.Result, n int,
	err error,
) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	result = examples.Message(str)
	return
}
