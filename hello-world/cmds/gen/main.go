package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
	structops "github.com/mus-format/musgen-go/options/struct"
	typeops "github.com/mus-format/musgen-go/options/type"
)

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.SayHelloCmd, cmds.SayFancyHelloCmd and the
// core.Cmd interface.
//
// For more details, see https://github.com/mus-format/musgen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/cmds"),
		genops.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// ValidateLength function will be used to validate the first Command field.
	// It protects the server from excessively large payloads - if
	// deserialization fails with an validation error, the corresponding client
	// connection will be closed.
	ops := structops.WithField(typeops.WithLenValidator("ValidateLength"))

	// Register core.Cmd interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		// Specify implementations.
		introps.WithStructImpl(reflect.TypeFor[cmds.SayHelloCmd](), ops),
		introps.WithStructImpl(reflect.TypeFor[cmds.SayFancyHelloCmd](), ops),
		// introps.WithRegisterMarshaller(), see the server-streaming example for
		// usage.
	)
	if err != nil {
		panic(err)
	}
	// g.RegisterInterface() method cannot be used in all cases — for example,
	// when a Command is already defined elsewhere. In such situations, use
	// g.AddInterface() instead (see the server-streaming example).

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
