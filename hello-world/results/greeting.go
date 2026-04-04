package results

// Greeting implements core.Result and exts.MarshallerTypedMUS interfaces.
type Greeting string

func (g Greeting) LastOne() bool {
	return true
}

func (g Greeting) String() string {
	return string(g)
}

// func (g Greeting) MarshalTypedMUS(w mus.Writer) (n int, err error) {
// 	return GreetingTypedMUS.Marshal(g, w)
// }

// func (g Greeting) SizeTypedMUS() (size int) {
// 	return GreetingTypedMUS.Size(g)
// }
