//go:generate go run gen/main.go
package cmds

import (
	com "github.com/mus-format/common-go"
)

const (
	TraceSayHelloCmdDTM com.DTM = iota + 10
	TraceSayFancyHelloCmdDTM
)
