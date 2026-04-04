package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/examples-go/hello-world/results"

	musgen "github.com/mus-format/mus-gen-go/mus"
	genopts "github.com/mus-format/mus-gen-go/options/gen"
	intropts "github.com/mus-format/mus-gen-go/options/interface"
)

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for results.Greeting and the core.Result interface.
//
// For more details, see https://github.com/mus-format/mus-gen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewGenerator(
		genopts.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/results"),
		genopts.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// Register core.Result interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Result](),
		// Specify implementations.
		intropts.WithDefinedTypeImpl(reflect.TypeFor[results.Greeting]()),
		// introps.WithRegisterMarshaller(), see the server-streaming example for
		// usage.
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
