package sqli

import (
	"fmt"
	"strings"
)

func AS(a StatementType, b StatementType) Statement {
	args := append(a.GetStatement().Args, b.GetStatement().Args...)
	return Statement{
		SQL:  fmt.Sprintf("%s AS %s", a.GetStatement().SQL, b.GetStatement().SQL),
		Args: args,
	}
}

func TABLE_WITH_COLUMNS(t TableType, columns ...ColumnType) Statement {
	columnsStr := ""
	for i, column := range columns {
		if len(columns) > 1 {
			if i == len(columns)-1 {
				columnsStr += column.GetNameAsAlias()
			} else {
				columnsStr += column.GetNameAsAlias() + ", "
			}
		} else {
			columnsStr += column.GetNameAsAlias()
		}
	}

	return Statement{
		SQL:  fmt.Sprintf("%s(%s)", t.GetName(), columnsStr),
		Args: []interface{}{},
	}
}

type ValuesSetSt []interface{}

func ValueSet(values ...interface{}) ValuesSetSt {
	return values
}

func VALUES(valueGroupList ...ValuesSetSt) Statement {
	valuesSql := []string{}
	valuesArgs := []interface{}{}

	for _, valueGroup := range valueGroupList {
		valueGroupSql := "("
		valueGroupArgs := []interface{}{}

		valueGroupLen := len(valueGroup)

		for i, value := range valueGroup {
			valueGroupSql += StringWithWithoutComma(valueGroupLen, i, QUERY_ARG)
			valueGroupArgs = append(valueGroupArgs, value)
		}

		valuesSql = append(valuesSql, valueGroupSql+")")
		valuesArgs = append(valuesArgs, valueGroupArgs...)
	}

	return Statement{
		SQL:  fmt.Sprintf("(VALUES %s)", strings.Join(valuesSql, ", ")),
		Args: valuesArgs,
	}
}

func VALUE[V any](c Column[V], value V) V {
	return value
}

func RETURNING(c ...ColumnType) Statement {
	sql := ""

	for i, column := range c {
		if len(c) > 1 {
			if i == len(c)-1 {
				sql += column.GetAlias()
			} else {
				sql += column.GetAlias() + ", "
			}
		} else {
			sql += column.GetAlias()
		}
	}

	stmt := fmt.Sprintf("RETURNING %s", sql)

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

func WITH(t Statement) Statement {
	return Statement{
		SQL:  fmt.Sprintf("WITH %s", t),
		Args: t.Args,
	}
}
