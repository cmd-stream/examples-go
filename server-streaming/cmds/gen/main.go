package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
)

// Main generates the mus-format.gen.go file with MUS serialization code for
// SayFancyHelloMultiCmd and the core.Cmd interface.
//
// For more details, see https://github.com/mus-format/musgen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/server-streaming/cmds"),
		genops.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genops.WithStream())
	if err != nil {
		panic(err)
	}

	// Register core.Cmd interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		// Specify implementations.
		introps.WithStructImpl(reflect.TypeFor[cmds.SayFancyHelloMultiCmd]()),
		introps.WithRegisterMarshaller(), // With this option all Commands must
		// implement the MarshallerTypedMUS interface from
		// github.com/mus-format/ext-stream-go.
		// It's not required and only affects how the Commands are serialized.
	)
	if err != nil {
		panic(err)
	}

	// Generate.
	bs, err := g.Generate()
	if err != nil {
		log.Println(string(bs))
	}

	// Write to file.
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	if err != nil {
		panic(err)
	}
}
