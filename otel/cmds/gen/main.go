// gen/main.go

package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	hwcmds "github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/otel/cmds"
	musgen "github.com/mus-format/mus-gen-go/mus"
	genopts "github.com/mus-format/mus-gen-go/options/gen"
	intropts "github.com/mus-format/mus-gen-go/options/interface"
)

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.TraceSayHelloCmd, cmds.TraceSayFancyHelloCmd,
// hwcmds.SayHelloCmd and core.Cmd interface.
//
// For more details, see https://github.com/mus-format/mus-gen-go.
func main() {
	// Create generator.
	g, err := musgen.NewGenerator(
		genopts.WithPkgPath("github.com/cmd-stream/examples-go/otel/cmds"),

		genopts.WithSerName(reflect.TypeFor[hwcmds.SayHelloCmd](),
			"hwcmds.SayHelloCmd"),
		genopts.WithSerName(reflect.TypeFor[hwcmds.SayFancyHelloCmd](),
			"hwcmds.SayFancyHelloCmd"),
		genopts.WithSerName(reflect.TypeFor[cmds.TraceSayHelloCmd](),
			"TraceSayHelloCmd"),
		genopts.WithSerName(reflect.TypeFor[cmds.TraceSayFancyHelloCmd](),
			"TraceSayFancyHelloCmd"),

		genopts.WithImportAlias("github.com/cmd-stream/examples-go/hello-world/cmds",
			"hwcmds"),
		genopts.WithImportAlias("github.com/cmd-stream/otelcmd-stream-go",
			"otelcmd"),

		genopts.WithStream(), // We're going to generate streaming code.
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

	err = g.AddTyped(traceSayHelloCmdType)
	if err != nil {
		panic(err)
	}

	// Add cmds.TraceSayFancyHelloCmd.
	traceSayFancyHelloCmdType := reflect.TypeFor[cmds.TraceSayFancyHelloCmd]()
	err = g.AddStruct(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddTyped(traceSayFancyHelloCmdType)
	if err != nil {
		panic(err)
	}

	// We can't use g.RegisterInterface() here because hwcmds.SayHelloCmd
	// has already been defined.
	//
	// This call instructs the generator to produce serializer for the core.Cmd
	// interface.
	err = g.AddInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		intropts.WithImpl(traceSayHelloCmdType),
		intropts.WithImpl(traceSayFancyHelloCmdType),
		intropts.WithImpl(reflect.TypeFor[hwcmds.SayHelloCmd]()),
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
	err = os.WriteFile("./mus.gen.go", bs, 0644)
	if err != nil {
		panic(err)
	}
}
