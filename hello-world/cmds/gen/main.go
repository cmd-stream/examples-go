package main

import (
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
	structops "github.com/mus-format/musgen-go/options/struct"
	typeops "github.com/mus-format/musgen-go/options/type"

	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// main function generates the mus-format.gen.go file, containing MUS
// serialization code for cmds.SayHelloCmd, cmds.SayFancyHelloCmd and the
// core.Cmd interface.
func main() {
	// Create a generator.
	g, err := musgen.NewCodeGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/hello-world/cmds"),
		genops.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genops.WithStream(), // We're going to generate streaming code.
	)
	assert.EqualError(err, nil)

	// ValidateLength function will be used to validate the first Command field.
	// It protects the server from excessively large payloads - if
	// deserialization fails with an validation error, the corresponding client
	// connection will be closed.
	ops := structops.WithField(typeops.WithLenValidator("ValidateLength"))

	// cmds.SayHelloCmd.
	sayHelloCmdType := reflect.TypeFor[cmds.SayHelloCmd]()
	err = g.AddStruct(sayHelloCmdType, ops)
	assert.EqualError(err, nil)

	// With this call the generator will produce SayHelloCmdDTS variable,
	// which helps to serialize 'DTM + SayHelloCmd'. DTS stands for Data Type
	// metadata Support.
	err = g.AddDTS(sayHelloCmdType)
	assert.EqualError(err, nil)

	// cmds.SayFancyHelloCmd.
	sayFancyHelloCmdType := reflect.TypeFor[cmds.SayFancyHelloCmd]()
	err = g.AddStruct(sayFancyHelloCmdType, ops)
	assert.EqualError(err, nil)

	err = g.AddDTS(sayFancyHelloCmdType)
	assert.EqualError(err, nil)

	// This call instructs the generator to produce serializer for the
	// core.Cmd interface.
	err = g.AddInterface(reflect.TypeFor[core.Cmd[receiver.Greeter]](),
		introps.WithImpl(sayHelloCmdType),
		introps.WithImpl(sayFancyHelloCmdType),
		introps.WithMarshaller(), /// SayHelloCmd and SayFancyHelloCmd should
		// also implement the MarshallerTypedMUS interface.
	)
	assert.EqualError(err, nil)

	// Generate.
	bs, err := g.Generate()
	assert.EqualError(err, nil)
	err = os.WriteFile("./mus-format.gen.go", bs, 0644)
	assert.EqualError(err, nil)
}
