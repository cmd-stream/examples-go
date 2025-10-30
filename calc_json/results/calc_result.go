package results

// CalcResult represents the Result of a calculation, returned from the server
// to the client. It implements the core.Result interface via its LastOne
// method.
type CalcResult int

func (r CalcResult) LastOne() bool { return true }
