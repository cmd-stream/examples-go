package receiver

import "strings"

// Greeter represents a Receiver and provides the functionality for
// creating greetings.
type Greeter interface {
	Interjection() string
	Adjective() string
	Join(strs ...string) string
}

func NewGreeter(interjection, adjective, sep string) greeter {
	return greeter{
		interjection: interjection,
		adjective:    adjective,
		sep:          sep,
	}
}

type greeter struct {
	interjection string
	adjective    string
	sep          string
}

func (g greeter) Interjection() string {
	return g.interjection
}

func (g greeter) Adjective() string {
	return g.adjective
}

func (g greeter) Join(strs ...string) string {
	return strings.Join(strs, g.sep)
}
