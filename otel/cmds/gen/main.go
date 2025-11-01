// gen/main.go

package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/otel/cmds"
	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
)

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.TraceSayHelloCmd, cmds.TraceSayFancyHelloCmd,
// hwcmds.SayHelloCmd and core.Cmd interface.
//
// For more details, see https://github.com/mus-format/musgen-go.
func main() {
	// Create generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/otel/cmds"),

		genops.WithSerName(reflect.TypeFor[hwcmds.SayHelloCmd](),
			"hwcmds.SayHelloCmd"),
		genops.WithSerName(reflect.TypeFor[hwcmds.SayFancyHelloCmd](),
			"hwcmds.SayFancyHelloCmd"),
		genops.WithSerName(reflect.TypeFor[cmds.TraceSayHelloCmd](),
			"TraceSayHelloCmd"),
		genops.WithSerName(reflect.TypeFor[cmds.TraceSayFancyHelloCmd](),
			"TraceSayFancyHelloCmd"),

		genops.WithImportAlias("github.com/cmd-stream/examples-go/hello-world/cmds",
			"hwcmds"),
		genops.WithImportAlias("github.com/cmd-stream/otelcmd-stream-go",
			"otelcmd"),

		genops.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// Add cmds.TraceSayHelloCmd.
	traceSayHelloCmdType := reflect.TypeFor[cmds.TraceSayHelloCmd]()
	err = g.AddStruct(traceSayHelloCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(traceSayHelloCmdType)
	if err != nil {
		panic(err)
	}

	// Add cmds.TraceSayFancyHelloCmd.
	traceSayFancyHelloCmdType := reflect.TypeFor[cmds.TraceSayFancyHelloCmd]()
	err = g.AddStruct(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	// We can't use g.RegisterInterface() here because hwcmds.SayHelloCmd
	// has already been defined.
	//
	// This call instructs the generator to produce serializer for the core.Cmd
	// interface.
	err = g.AddInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		introps.WithImpl(traceSayHelloCmdType),
		introps.WithImpl(traceSayFancyHelloCmdType),
		introps.WithImpl(reflect.TypeFor[hwcmds.SayHelloCmd]()),
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
	if err != nil {
		panic(err)
	}
}
