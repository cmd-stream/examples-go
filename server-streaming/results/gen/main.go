package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/examples-go/server-streaming/results"

	musgen "github.com/mus-format/mus-gen-go/mus"
	genops "github.com/mus-format/mus-gen-go/options/gen"
	introps "github.com/mus-format/mus-gen-go/options/interface"
)

// Main generates the mus-format.gen.go file with MUS serialization code for
// Greeting and the core.Result interface.
//
// For more details, see https://github.com/mus-format/mus-gen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/server-streaming/results"),
		genops.WithStream())
	if err != nil {
		panic(err)
	}

	// Register core.Result interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Result](),
		// Specify implementations.
		introps.WithStructImpl(reflect.TypeFor[results.Greeting]()),
		introps.WithRegisterMarshaller(), // With this option all Results must
		// implement the MarshallerTypedMUS interface from
		// github.com/mus-format/ext-stream-go.
		// It's not required and only affects how the Results are serialized.
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
