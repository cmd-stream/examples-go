//go:generate go run gen/main.go
package cmds

import com "github.com/mus-format/common-go"

const (
	AddCmdDTM com.DTM = iota + 1
	SubCmdDTM
)
