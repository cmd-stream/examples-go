package cmds

import "errors"

// ValidateLength is used on the server to validate a Command's content.
func ValidateLength(length int) (err error) {
	if length > 1000 {
		return errors.New("too large")
	}
	return
}
