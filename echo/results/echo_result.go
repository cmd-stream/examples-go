package results

type EchoResult string

func (r EchoResult) LastOne() bool {
	return true
}
