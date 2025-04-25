package sqlification

import "fmt"

func UPDATE(t TableType) Statement {
	return Statement{
		SQL:  fmt.Sprintf("UPDATE %s", t.GetName()),
		Args: []interface{}{},
	}
}

func SET(c ...Statement) Statement {
	result := Statement{
		SQL:  "SET ",
		Args: []interface{}{},
	}

	for i, query := range c {
		if i == 0 {
			result.SQL += query.SQL
		} else {
			result.SQL += (", " + query.SQL)
		}
		result.Args = append(result.Args, query.Args...)
	}

	return result
}

func SET_VALUE[V any](c Column[V], value V) Statement {
	return Statement{
		SQL:  fmt.Sprintf("%s = %s", c.GetName(), QUERY_ARG),
		Args: []interface{}{value},
	}
}

func SET_VALUE_COLUMN[A any, B any](a Column[A], b Column[B]) Statement {
	return Statement{
		SQL:  fmt.Sprintf("%s = %s.%s", a.GetName(), b.TableAlias, b.GetName()),
		Args: []interface{}{},
	}
}
