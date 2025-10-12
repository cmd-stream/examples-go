package main

import (
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
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/results"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// results.Greeting.
	greetingType := reflect.TypeFor[results.Greeting]()
	err = g.AddDefinedType(greetingType)
	assert.EqualError(err, nil)

	err = g.AddDTS(greetingType)
	assert.EqualError(err, nil)

	// core.Result interface.
	err = g.AddInterface(reflect.TypeFor[core.Result](),
		introps.WithImpl(greetingType),
		introps.WithMarshaller(),
	)
	assert.EqualError(err, nil)

	// Generate.
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
