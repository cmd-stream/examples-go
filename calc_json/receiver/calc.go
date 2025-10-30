package receiver

// Calc serves as the application's "Receiver" or service layer.
//
// It encapsulates the core application logic, ensuring it remains completely
// decoupled from Command structures, transport, and network concerns.
// The Calc instance is dependency-injected into Command.Exec methods,
// which maximizes reusability and simplifies unit testing.
type Calc interface {
	// Add performs the addition.
	Add(a, b int) int

	// Sub performs the subtraction.
	Sub(a, b int) int
}

func NewCalc() calc {
	return calc{}
}

type calc struct{}

func (c calc) Add(a, b int) int { return a + b }

func (c calc) Sub(a, b int) int { return a - b }
