package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/calc/results"
	musgen "github.com/mus-format/musgen-go/mus"

	introps "github.com/mus-format/musgen-go/options/interface"

	genops "github.com/mus-format/musgen-go/options/generate"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for results.Result and core.Result interface.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/calc/results"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// results.CalcResult
	calcResultType := reflect.TypeFor[results.CalcResult]()
	err = g.AddDefinedType(calcResultType)
	assert.EqualError(err, nil)

	// With this call the generator will produce ResultDTS variable,
	// which helps to serialize 'DTM + CalcResult'. DTS stands for Data Type
	// metadata Support.
	err = g.AddDTS(calcResultType)
	assert.EqualError(err, nil)
	//
	// This call instructs the generator to produce serializer for the
	// core.Cmd interface.
	err = g.AddInterface(reflect.TypeFor[core.Result](),
		introps.WithImpl(calcResultType),
	)
	assert.EqualError(err, nil)

	// Generate
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
