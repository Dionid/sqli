package sqlification

import (
	"fmt"
	"strings"
)

var QUERY_ARG = "$$DIO_ARG$$"

func StringWithWithoutComma(len int, i int, val string) string {
	if len == 1 {
		return val
	}

	if i == len-1 {
		return val
	}

	return val + ", "
}

func TrimQuotes(val string) string {
	return strings.Trim(val, `"`)
}

func DoubleQuotes(val string) string {
	if val == "*" {
		return val
	}
	return fmt.Sprintf(`"%s"`, TrimQuotes(val))
}
