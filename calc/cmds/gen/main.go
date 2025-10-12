package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/calc/cmds"
	rcvr "github.com/cmd-stream/examples-go/calc/receiver"
	musgen "github.com/mus-format/musgen-go/mus"

	introps "github.com/mus-format/musgen-go/options/interface"

	genops "github.com/mus-format/musgen-go/options/generate"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.AddCmd, cmds.SubCmd Commands and core.Cmd
// interface.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/calc/cmds"),
		genops.WithImport("github.com/cmd-stream/examples-go/calc/receiver"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// cmds.AddCmd
	addCmdType := reflect.TypeFor[cmds.AddCmd]()
	err = g.AddStruct(addCmdType)
	assert.EqualError(err, nil)

	// With this call the generator will produce AddCmdDTS variable,
	// which helps to serialize 'DTM + AddCmd'. DTS stands for Data Type
	// metadata Support.
	err = g.AddDTS(addCmdType)
	assert.EqualError(err, nil)

	// cmds.SubCmd
	subCmdType := reflect.TypeFor[cmds.SubCmd]()
	err = g.AddStruct(subCmdType)
	assert.EqualError(err, nil)

	err = g.AddDTS(subCmdType)
	assert.EqualError(err, nil)

	// This call instructs the generator to produce serializer for the
	// core.Cmd interface.
	err = g.AddInterface(reflect.TypeFor[core.Cmd[rcvr.Calc]](),
		introps.WithImpl(addCmdType),
		introps.WithImpl(subCmdType),
	)
	assert.EqualError(err, nil)

	// Generate
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
