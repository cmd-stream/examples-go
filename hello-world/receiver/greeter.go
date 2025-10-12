package receiver

import "strings"

// Greeter serves as the application's "Receiver" or application layer.
//
// It holds the core application logic, decoupled from Commands and network
// transport. Defining it as an interface ensures both the Greeter
// implementation and its dependent Commands can be easily tested in isolation.
type Greeter interface {
	// Interjection returns the fixed interjection (e.g., "Hi").
	Interjection() string
	// Adjective returns the fixed adjective (e.g., "amazing").
	Adjective() string
	// Join concatenates strings using the configured separator.
	Join(strs ...string) string
}

// NewGreeter returns a new greeter.
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
