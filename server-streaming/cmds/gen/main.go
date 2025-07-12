package main

import (
	"log"
	"os"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/server-streaming/cmds"

	musgen "github.com/mus-format/musgen-go/mus"
	genops "github.com/mus-format/musgen-go/options/generate"
	introps "github.com/mus-format/musgen-go/options/interface"
)

// The main function generates the mus-format.gen.go file containing MUS
// serialization code for SayHelloCmd, SayFancyHelloCmd, and Result.
func main() {
	// Create a generator.
	g, err := musgen.NewFileGenerator(
		genops.WithPkgPath("github.com/cmd-stream/examples-go/server-streaming/cmds"),
		genops.WithImport("github.com/cmd-stream/examples-go/hello-world/receiver"),
		genops.WithStream())
	if err != nil {
		panic(err)
	}

	// Add SayFancyHelloMultiCmd.
	sayFancyHelloMultiCmdType := reflect.TypeFor[cmds.SayFancyHelloMultiCmd]()
	err = g.AddStruct(sayFancyHelloMultiCmdType)
	if err != nil {
		panic(err)
	}

	err = g.AddDTS(sayFancyHelloMultiCmdType)
	if err != nil {
		panic(err)
	}

	// Add core.Cmd
	err = g.AddInterface(reflect.TypeFor[core.Cmd[receiver.Greeter]](),
		introps.WithImpl(sayFancyHelloMultiCmdType))
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
