package results

import "github.com/mus-format/mus-stream-go"

// NewGreeting creates a new Greeting.
func NewGreeting(str string, lastOne bool) Greeting {
	return Greeting{str, lastOne}
}

// Greeting implements core.Result and ext.MarshallerTypedMUS interfaces.
//
// We have to define MarshalTypedMUS and SizeTypedMUS methods (implement the
// MarshallerTypedMUS interface) because the core.Result interface
// serialization code was generated with introps.WithRegisterMarshaller().
type Greeting struct {
	str     string
	lastOne bool
}

func (g Greeting) String() string {
	return g.str
}

func (g Greeting) LastOne() bool {
	return g.lastOne
}

func (g Greeting) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return GreetingDTS.Marshal(g, w)
}

func (g Greeting) SizeTypedMUS() (size int) {
	return GreetingDTS.Size(g)
}
