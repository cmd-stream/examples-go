package receiver

import "strings"

// Greeter is the application's Receiver. It defines the core logic,
// kept independent from Commands and network transport.
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
