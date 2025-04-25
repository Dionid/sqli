package sqli

import "fmt"

func COUNT[V any](column Column[V]) Statement {
	stmt := fmt.Sprintf("COUNT(%s)", column.GetAliasWithTableAlias())

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}
