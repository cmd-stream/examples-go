package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"

	musgen "github.com/mus-format/mus-gen-go/mus"
	genopts "github.com/mus-format/mus-gen-go/options/gen"
	intropts "github.com/mus-format/mus-gen-go/options/interface"
)

// Main generates the mus-format.gen.go file with MUS serialization code for
// SayFancyHelloMultiCmd and the core.Cmd interface.
//
// For more details, see https://github.com/mus-format/mus-gen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewGenerator(
		genopts.WithPkgPath("github.com/cmd-stream/examples-go/server-streaming/cmds"),
		genopts.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genopts.WithStream())
	if err != nil {
		panic(err)
	}

	// Register core.Cmd interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		// Specify implementations.
		intropts.WithStructImpl(reflect.TypeFor[cmds.SayFancyHelloMultiCmd]()),
		intropts.WithRegisterMarshaller(), // With this option all Commands must
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
	err = os.WriteFile("./mus.gen.go", bs, 0644)
	if err != nil {
		panic(err)
	}
}
