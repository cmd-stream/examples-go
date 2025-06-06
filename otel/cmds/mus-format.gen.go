// Code generated by musgen-go. DO NOT EDIT.

package cmds

import (
	"fmt"

	"github.com/cmd-stream/core-go"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	otelcmd "github.com/cmd-stream/otelcmd-stream-go"
	dts "github.com/mus-format/dts-stream-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

var (
	map98ΣX2OmXAmD4Cz9Lor4fZwΞΞ = ord.NewMapSer[string, string](ord.String, ord.String)
	ptrFjtJRFat1jHEPoFCkftt6gΞΞ = ord.NewPtrSer[map[string]string](map98ΣX2OmXAmD4Cz9Lor4fZwΞΞ)
)

var TraceSayHelloCmdMUS = traceSayHelloCmdMUS{}

type traceSayHelloCmdMUS struct{}

func (s traceSayHelloCmdMUS) Marshal(v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd], w muss.Writer) (n int, err error) {
	n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Marshal(v.MapCarrier, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = hwcmds.SayHelloCmdMUS.Marshal(v.Cmd, w)
	n += n1
	return
}

func (s traceSayHelloCmdMUS) Unmarshal(r muss.Reader) (v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd], n int, err error) {
	v.MapCarrier, n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	v.Cmd, n1, err = hwcmds.SayHelloCmdMUS.Unmarshal(r)
	n += n1
	return
}

func (s traceSayHelloCmdMUS) Size(v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd]) (size int) {
	size = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Size(v.MapCarrier)
	return size + hwcmds.SayHelloCmdMUS.Size(v.Cmd)
}

func (s traceSayHelloCmdMUS) Skip(r muss.Reader) (n int, err error) {
	n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Skip(r)
	if err != nil {
		return
	}
	var n1 int
	n1, err = hwcmds.SayHelloCmdMUS.Skip(r)
	n += n1
	return
}

var TraceSayHelloCmdDTS = dts.New[otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd]](TraceSayHelloCmdDTM, TraceSayHelloCmdMUS)

var TraceSayFancyHelloCmdMUS = traceSayFancyHelloCmdMUS{}

type traceSayFancyHelloCmdMUS struct{}

func (s traceSayFancyHelloCmdMUS) Marshal(v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd], w muss.Writer) (n int, err error) {
	n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Marshal(v.MapCarrier, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = hwcmds.SayFancyHelloCmdMUS.Marshal(v.Cmd, w)
	n += n1
	return
}

func (s traceSayFancyHelloCmdMUS) Unmarshal(r muss.Reader) (v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd], n int, err error) {
	v.MapCarrier, n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	v.Cmd, n1, err = hwcmds.SayFancyHelloCmdMUS.Unmarshal(r)
	n += n1
	return
}

func (s traceSayFancyHelloCmdMUS) Size(v otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd]) (size int) {
	size = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Size(v.MapCarrier)
	return size + hwcmds.SayFancyHelloCmdMUS.Size(v.Cmd)
}

func (s traceSayFancyHelloCmdMUS) Skip(r muss.Reader) (n int, err error) {
	n, err = ptrFjtJRFat1jHEPoFCkftt6gΞΞ.Skip(r)
	if err != nil {
		return
	}
	var n1 int
	n1, err = hwcmds.SayFancyHelloCmdMUS.Skip(r)
	n += n1
	return
}

var TraceSayFancyHelloCmdDTS = dts.New[otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd]](TraceSayFancyHelloCmdDTM, TraceSayFancyHelloCmdMUS)

var CmdMUS = cmdMUS{}

type cmdMUS struct{}

func (s cmdMUS) Marshal(v core.Cmd[receiver.Greeter], w muss.Writer) (n int, err error) {
	switch t := v.(type) {
	case otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd]:
		return TraceSayHelloCmdDTS.Marshal(t, w)
	case otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd]:
		return TraceSayFancyHelloCmdDTS.Marshal(t, w)
	case hwcmds.SayHelloCmd:
		return hwcmds.SayHelloCmdDTS.Marshal(t, w)
	default:
		panic(fmt.Sprintf("unexpected %v type", t))
	}
}

func (s cmdMUS) Unmarshal(r muss.Reader) (v core.Cmd[receiver.Greeter], n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case TraceSayHelloCmdDTM:
		v, n1, err = TraceSayHelloCmdDTS.UnmarshalData(r)
	case TraceSayFancyHelloCmdDTM:
		v, n1, err = TraceSayFancyHelloCmdDTS.UnmarshalData(r)
	case hwcmds.SayHelloCmdDTM:
		v, n1, err = hwcmds.SayHelloCmdDTS.UnmarshalData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}

func (s cmdMUS) Size(v core.Cmd[receiver.Greeter]) (size int) {
	switch t := v.(type) {
	case otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayHelloCmd]:
		return TraceSayHelloCmdDTS.Size(t)
	case otelcmd.TraceCmd[receiver.Greeter, hwcmds.SayFancyHelloCmd]:
		return TraceSayFancyHelloCmdDTS.Size(t)
	case hwcmds.SayHelloCmd:
		return hwcmds.SayHelloCmdDTS.Size(t)
	default:
		panic(fmt.Sprintf("unexpected %v type", t))
	}
}

func (s cmdMUS) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case TraceSayHelloCmdDTM:
		n1, err = TraceSayHelloCmdDTS.SkipData(r)
	case TraceSayFancyHelloCmdDTM:
		n1, err = TraceSayFancyHelloCmdDTS.SkipData(r)
	case hwcmds.SayHelloCmdDTM:
		n1, err = hwcmds.SayHelloCmdDTS.SkipData(r)
	default:
		err = fmt.Errorf("unexpected %v DTM", dtm)
		return
	}
	n += n1
	return
}
