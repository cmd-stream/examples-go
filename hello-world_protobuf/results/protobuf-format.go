package results

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
	hwresults "github.com/cmd-stream/examples-go/hello-world/results"
)

// Greeting

var GreetingProtobuf = greetingProtobuf{}

type greetingProtobuf struct{}

func (s greetingProtobuf) Marshal(result Greeting, w muss.Writer) (n int,
	err error) {
	bs, err := proto.Marshal(result.GreetingData)
	if err != nil {
		return
	}
	l := len(bs)
	n, err = varint.PositiveInt.Marshal(l, w)
	if err != nil {
		return
	}
	n1, err := w.Write(bs)
	n += n1
	return
}

func (s greetingProtobuf) Unmarshal(r muss.Reader) (result Greeting, n int,
	err error) {
	data := &GreetingData{}
	n, err = unmarshal[*GreetingData](data, nil, r)
	if err != nil {
		return
	}
	result.GreetingData = data
	return
}

func (s greetingProtobuf) Size(result Greeting) (size int) {
	panic("not implemented")
}

func (s greetingProtobuf) Skip(r muss.Reader) (n int, err error) {
	panic("not implemented")
}

var GreetingDTS = dts.New[Greeting](hwresults.GreetingDTM, GreetingProtobuf)

// core.Result

var ResultProtobuf = resultProtobuf{}

type resultProtobuf struct{}

func (s resultProtobuf) Marshal(v core.Result, w muss.Writer) (n int, err error) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.MarshalTypedProtobuf(w)
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s resultProtobuf) Unmarshal(r muss.Reader) (v core.Result, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hwresults.GreetingDTM:
		v, n1, err = GreetingDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s resultProtobuf) Size(v core.Result) (size int) {
	if m, ok := v.(exts.MarshallerTypedProtobuf); ok {
		return m.SizeTypedProtobuf()
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedProtobuf interface",
		reflect.TypeOf(v)))
}

func (s resultProtobuf) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case hwresults.GreetingDTM:
		n1, err = GreetingDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
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
