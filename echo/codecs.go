package main

import (
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w transport.Writer) (
	n int, err error,
) {
	_, err = ord.String.Marshal(string(result.(Message)), w)
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd core.Cmd[struct{}], n int,
	err error,
) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	cmd = Message(str)
	return
}

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd core.Cmd[struct{}], w transport.Writer) (
	n int, err error,
) {
	_, err = ord.String.Marshal(string(cmd.(Message)), w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result core.Result, n int,
	err error,
) {
	str, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	result = Message(str)
	return
}
