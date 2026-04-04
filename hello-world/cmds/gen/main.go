package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	rcvr "github.com/cmd-stream/examples-go/hello-world/receiver"

	musgen "github.com/mus-format/mus-gen-go/mus"
	fdopts "github.com/mus-format/mus-gen-go/options/field"
	genopts "github.com/mus-format/mus-gen-go/options/gen"
	intropts "github.com/mus-format/mus-gen-go/options/interface"
	stopts "github.com/mus-format/mus-gen-go/options/struct"
	tpopts "github.com/mus-format/mus-gen-go/options/type"
)

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.SayHelloCmd, cmds.SayFancyHelloCmd and the
// core.Cmd interface.
//
// For more details, see https://github.com/mus-format/mus-gen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewGenerator(
		genopts.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/cmds"),
		genopts.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genopts.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// ValidateLength function will be used to validate the first Command field.
	// It protects the server from excessively large payloads - if
	// deserialization fails with an validation error, the corresponding client
	// connection will be closed.
	opts := stopts.WithField(
		fdopts.WithType(
			tpopts.WithLenValidator("ValidateLength"),
		),
	)

	// Register core.Cmd interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Cmd[rcvr.Greeter]](),
		// Specify implementations.
		intropts.WithStructImpl(reflect.TypeFor[cmds.SayHelloCmd](), opts),
		intropts.WithStructImpl(reflect.TypeFor[cmds.SayFancyHelloCmd](), opts),
		// introps.WithRegisterMarshaller(), see the server-streaming example for
		// usage.
	)
	if err != nil {
		panic(err)
	}
	// g.RegisterInterface() method cannot be used in all cases — for example,
	// when a Command is already defined elsewhere. In such situations, use
	// g.AddInterface() instead (see the otel example).

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
