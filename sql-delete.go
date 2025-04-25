package sqlification

import "fmt"

// DELETE_FROM
func DELETE_FROM(t TableType) Statement {
	return Statement{
		SQL:  fmt.Sprintf("DELETE FROM %s", t.GetName()),
		Args: []interface{}{},
	}
}
