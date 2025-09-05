package results

import muss "github.com/mus-format/mus-stream-go"

// NewGreeting creates a new Greeting.
func NewGreeting(str string) Greeting {
	return Greeting{
		GreetingData: &GreetingData{Str: str},
	}
}

// Greeting implements core.Result and exts.MarshallerTypedProtobuf interfaces.
type Greeting struct {
	*GreetingData
}

func (g Greeting) String() string {
	return g.Str
}

func (g Greeting) LastOne() bool {
	return true
}

func (g Greeting) MarshalTypedProtobuf(w muss.Writer) (n int, err error) {
	return GreetingDTS.Marshal(g, w)
}

func (g Greeting) SizeTypedProtobuf() (size int) {
	return GreetingDTS.Size(g)
}
