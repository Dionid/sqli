package sqli

import (
	"fmt"
)

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("(validation) %s", e.Msg)
}
