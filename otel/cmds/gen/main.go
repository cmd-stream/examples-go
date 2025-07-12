// gen/main.go

package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/otel/cmds"
	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"

	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for hw.SayHelloCmd, hw.SayFancyHelloCmd, hw.Greeting,
// core.Cmd and core.Result interfaces.
func main() {
	// Create generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/otel/cmds"),

		genops.WithSerName(reflect.TypeFor[hwcmds.SayHelloCmd](), "hwcmds.SayHelloCmd"),
		genops.WithSerName(reflect.TypeFor[hwcmds.SayFancyHelloCmd](), "hwcmds.SayFancyHelloCmd"),
		genops.WithSerName(reflect.TypeFor[cmds.TraceSayHelloCmd](), "TraceSayHelloCmd"),
		genops.WithSerName(reflect.TypeFor[cmds.TraceSayFancyHelloCmd](), "TraceSayFancyHelloCmd"),

		genops.WithImportAlias("github.com/cmd-stream/examples-go/hello-world/cmds", "hwcmds"),
		genops.WithImportAlias("github.com/cmd-stream/otelcmd-stream-go", "otelcmd"),

		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// Add hw.SayHelloCmd.
	traceSayHelloCmdType := reflect.TypeFor[cmds.TraceSayHelloCmd]()
	err = g.AddStruct(traceSayHelloCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(traceSayHelloCmdType)
	if err != nil {
		panic(err)
	}

	// Add hw.SayHelloCmd.
	traceSayFancyHelloCmdType := reflect.TypeFor[cmds.TraceSayFancyHelloCmd]()
	err = g.AddStruct(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	// Add core.Cmd. This call instructs the generator to produce serializer for
	// the core.Cmd interface.
	err = g.AddInterface(reflect.TypeFor[core.Cmd[receiver.Greeter]](),
		introps.WithImpl(traceSayHelloCmdType),
		introps.WithImpl(traceSayFancyHelloCmdType),
		introps.WithImpl(reflect.TypeFor[hwcmds.SayHelloCmd]()),
		// introps.WithMarshaller(), // SayHelloCmd and SayFancyHelloCmd should
		// also implement the MarshallerTypedMUS interface. More on this later.
	)
	if err != nil {
		panic(err)
	}

	// Generate.
	bs, err := g.Generate()
	if err != nil {
		log.Println(err)
	}

	// Write to file.
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
