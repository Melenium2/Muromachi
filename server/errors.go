package server

import "fmt"

var (
	ErrInternalError = fmt.Errorf("%s", "we have some internal troubles, sorry for that")
)
