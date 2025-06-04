package results

import muss "github.com/mus-format/mus-stream-go"

// NewGreeting creates a new Greeting.
func NewGreeting(str string, lastOne bool) Greeting {
	return Greeting{str, lastOne}
}

// Greeting implements the Greeting interface.
type Greeting struct {
	str     string
	lastOne bool
}

func (r Greeting) String() string {
	return r.str
}

// Command in this tutorial sends back several results.
func (r Greeting) LastOne() bool {
	return r.lastOne
}

func (c Greeting) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return GreetingDTS.Marshal(c, w)
}

func (c Greeting) SizeTypedMUS() (size int) {
	return GreetingDTS.Size(c)
}
