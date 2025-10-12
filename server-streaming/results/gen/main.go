package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/server-streaming/results"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
)

// Main generates the mus-format.gen.go file with MUS serialization code for
// Greeting and the core.Result interface.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/server-streaming/results"),
		genops.WithStream())
	if err != nil {
		panic(err)
	}

	// Add Greeting.
	greetingType := reflect.TypeFor[results.Greeting]()
	err = g.AddStruct(greetingType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(greetingType)
	if err != nil {
		panic(err)
	}

	// Add core.Result.
	err = g.AddInterface(reflect.TypeFor[core.Result](),
		introps.WithImpl(greetingType),
		introps.WithMarshaller(),
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
