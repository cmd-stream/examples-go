package cmds

import (
	"fmt"
	"io"
	reflect "reflect"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-stream-go"
	exts "github.com/mus-format/ext-protobuf-stream-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
	"google.golang.org/protobuf/proto"

	"github.com/cmd-stream/core-go"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	receiver "github.com/cmd-stream/examples-go/hello-world/receiver"
)

// SayHelloCmd

var SayHelloCmdProtobuf = sayHelloCmdProtobuf{}

type sayHelloCmdProtobuf struct{}

func (s sayHelloCmdProtobuf) Marshal(c SayHelloCmd, w muss.Writer) (n int,
	err error) {
	return marshal(c, w)
}

func (s sayHelloCmdProtobuf) Unmarshal(r muss.Reader) (c SayHelloCmd, n int,
	err error) {
	data := &SayHelloData{}
	n, err = unmarshal(data, com.ValidatorFn[int](hwcmds.ValidateLength), r)
	if err != nil {
		return
	}
	c.SayHelloData = data
	return
}

func (s sayHelloCmdProtobuf) Size(c SayHelloCmd) (size int) {
	panic("not implemented")
}

func (s sayHelloCmdProtobuf) Skip(r muss.Reader) (n int,
	err error) {
	panic("not implemented")
}

// SayFancyHelloCmd

var SayFancyHelloCmdProtobuf = sayFancyHelloCmdProtobuf{}

type sayFancyHelloCmdProtobuf struct{}

func (s sayFancyHelloCmdProtobuf) Marshal(c SayFancyHelloCmd, w muss.Writer) (
	n int, err error) {
	return marshal(c, w)
}

func (s sayFancyHelloCmdProtobuf) Unmarshal(r muss.Reader) (c SayFancyHelloCmd,
	n int, err error) {
	data := &SayFancyHelloData{}
	n, err = unmarshal(data, com.ValidatorFn[int](hwcmds.ValidateLength), r)
	if err != nil {
		return
	}
	c.SayFancyHelloData = data
	return
}

func (s sayFancyHelloCmdProtobuf) Size(c SayFancyHelloCmd) (size int) {
	panic("not implemented")
}

func (s sayFancyHelloCmdProtobuf) Skip(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

// DTS

var (
	SayHelloCmdDTS = dts.New[SayHelloCmd](hwcmds.SayHelloCmdDTM,
		SayHelloCmdProtobuf)
	SayFancyHelloCmdDTS = dts.New[SayFancyHelloCmd](hwcmds.SayFancyHelloCmdDTM,
		SayFancyHelloCmdProtobuf)
)

// core.Cmd

var CmdProtobuf = cmdProtobuf{}

type cmdProtobuf struct{}

func (s cmdProtobuf) Marshal(v core.Cmd[receiver.Greeter], w muss.Writer) (n int,
	err error) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.MarshalTypedProtobuf(w)
	}
	str := fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v))
	panic(str)
}

func (s cmdProtobuf) Unmarshal(r muss.Reader) (v core.Cmd[receiver.Greeter], n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hwcmds.SayHelloCmdDTM:
		v, n1, err = SayHelloCmdDTS.UnmarshalData(r)
	case hwcmds.SayFancyHelloCmdDTM:
		v, n1, err = SayFancyHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s cmdProtobuf) Size(v core.Cmd[receiver.Greeter]) (size int) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.SizeTypedProtobuf()
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s cmdProtobuf) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hwcmds.SayHelloCmdDTM:
		n1, err = SayHelloCmdDTS.SkipData(r)
	case hwcmds.SayFancyHelloCmdDTM:
		n1, err = SayFancyHelloCmdDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func marshal[T proto.Message](c T, w muss.Writer) (n int, err error) {
	bs, err := proto.Marshal(c)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.PositiveInt.Marshal(l, w)
	if err != nil {
		return
	}
	_, err = w.Write(bs)
	n += l
	return
}

func unmarshal[T proto.Message](d T, v com.Validator[int],
	r muss.Reader) (n int, err error) {
	l, n, err := varint.PositiveInt.Unmarshal(r)
	if err != nil {
		return
	}
	if v != nil {
		if err = v.Validate(l); err != nil {
			return
		}
	}
	bs := make([]byte, l)
	n1, err := io.ReadFull(r, bs)
	n += n1
	if err != nil {
		return
	}
	err = proto.Unmarshal(bs, d)
	return
}
