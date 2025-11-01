package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/results"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"

	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for results.Greeting and the core.Result interface.
//
// For more details, see https://github.com/mus-format/musgen-go.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/results"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	if err != nil {
		panic(err)
	}

	// Register core.Result interface.
	err = g.RegisterInterface(reflect.TypeFor[core.Result](),
		// Specify implementations.
		introps.WithDefinedTypeImpl(reflect.TypeFor[results.Greeting]()),
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
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	if err != nil {
		panic(err)
	}
}
