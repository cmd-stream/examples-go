//go:generate go run gen/main.go
package cmds

import (
	"time"

	com "github.com/mus-format/common-go"
)

const (
	SayHelloCmdDTM com.DTM = iota + 1
	SayFancyHelloCmdDTM
)

// CmdExecDuration defines the duration of Command execution.
const CmdExecDuration = time.Second
