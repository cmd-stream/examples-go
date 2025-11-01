//go:generate go run gen/main.go
package cmds

import (
	"errors"
	"time"
)

// CmdExecDuration defines the duration of Command execution.
const CmdExecDuration = time.Second

// ValidateLength is used on the server to validate a Command's content.
func ValidateLength(length int) (err error) {
	if length > 1000 {
		return errors.New("too large")
	}
	return
}
