package results

import muss "github.com/mus-format/mus-stream-go"

// NewGreeting creates a new Greeting.
func NewGreeting(str string, lastOne bool) Greeting {
	return Greeting{str, lastOne}
}

// Greeting implements core.Result and exts.MarshallerTypedMUS interfaces.
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

func (g Greeting) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return GreetingDTS.Marshal(g, w)
}

func (g Greeting) SizeTypedMUS() (size int) {
	return GreetingDTS.Size(g)
}
