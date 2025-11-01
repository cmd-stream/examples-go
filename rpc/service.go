package main

import (
	"context"

	"github.com/cmd-stream/examples-go/hello-world/cmds"
	"github.com/cmd-stream/examples-go/hello-world/receiver"
	"github.com/cmd-stream/examples-go/hello-world/results"

	sndr "github.com/cmd-stream/sender-go"
)

type GreeterService struct {
	sender sndr.Sender[receiver.Greeter]
}

func (s GreeterService) SayHello(ctx context.Context, str string) (string, error) {
	cmd := cmds.SayHelloCmd{Str: str}
	greeting, err := s.sender.Send(ctx, cmd)
	if err != nil {
		return "", err
	}
	return greeting.(results.Greeting).String(), nil
}
