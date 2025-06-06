// Code generated by musgen-go. DO NOT EDIT.

package results

import (
	"fmt"
	"reflect"

	"github.com/cmd-stream/core-go"
	dts "github.com/mus-format/dts-stream-go"
	exts "github.com/mus-format/ext-mus-stream-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

var GreetingMUS = greetingMUS{}

type greetingMUS struct{}

func (s greetingMUS) Marshal(v Greeting, w muss.Writer) (n int, err error) {
	return ord.String.Marshal(string(v), w)
}

func (s greetingMUS) Unmarshal(r muss.Reader) (v Greeting, n int, err error) {
	tmp, n, err := ord.String.Unmarshal(r)
	if err != nil {
		return
	}
	v = Greeting(tmp)
	return
}

func (s greetingMUS) Size(v Greeting) (size int) {
	return ord.String.Size(string(v))
}

func (s greetingMUS) Skip(r muss.Reader) (n int, err error) {
	return ord.String.Skip(r)
}

var GreetingDTS = dts.New[Greeting](GreetingDTM, GreetingMUS)

var ResultMUS = resultMUS{}

type resultMUS struct{}

func (s resultMUS) Marshal(v core.Result, w muss.Writer) (n int, err error) {
	if m, ok := v.(exts.MarshallerTypedMUS); ok {
		return m.MarshalTypedMUS(w)
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedMUS interface", reflect.TypeOf(v)))
}

func (s resultMUS) Unmarshal(r muss.Reader) (v core.Result, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case GreetingDTM:
		v, n1, err = GreetingDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s resultMUS) Size(v core.Result) (size int) {
	if m, ok := v.(exts.MarshallerTypedMUS); ok {
		return m.SizeTypedMUS()
	}
	panic(fmt.Sprintf("%v doesn't implement the exts.MarshallerTypedMUS interface", reflect.TypeOf(v)))
}

func (s resultMUS) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case GreetingDTM:
		n1, err = GreetingDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}
