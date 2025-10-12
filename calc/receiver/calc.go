package receiver

// Calc serves as the application's "Receiver" or service layer.
//
// It encapsulates the core application logic, ensuring it remains completely
// decoupled from Command structures, transport, and network concerns.
// The Calc instance is dependency-injected into Command.Exec methods,
// which maximizes reusability and simplifies unit testing.
type Calc struct{}

// Add performs the addition.
func (c Calc) Add(a, b int) int { return a + b }

// Sub performs the subtraction.
func (c Calc) Sub(a, b int) int { return a - b }
