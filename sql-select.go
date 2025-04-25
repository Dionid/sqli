package sqli

import (
	"fmt"
)

// # SELECT

func SELECT(columns ...ColumnType) Statement {
	columnsStr := ""

	for i, column := range columns {
		if len(columns) > 1 {
			if i == len(columns)-1 {
				columnsStr += column.GetAliasWithTableAlias()
			} else {
				columnsStr += column.GetAliasWithTableAlias() + ", "
			}
		} else {
			columnsStr += column.GetAliasWithTableAlias()
		}
	}

	stmt := fmt.Sprintf("SELECT %s", columnsStr)

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

// # FROM

func FROM(table NameWithAliaser) Statement {
	stmt := fmt.Sprintf("FROM %s", table.GetNameWithAlias())

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

func FROM_RAW(sq Statement) Statement {
	stmt := fmt.Sprintf("FROM %s", sq.SQL)

	return Statement{
		SQL:  stmt,
		Args: sq.Args,
	}
}

// ORDER BY

type OrderDirection struct {
	Value string
}

var ASC = OrderDirection{
	Value: "ASC",
}

var DESC = OrderDirection{
	Value: "DESC",
}

type ColumnOrder struct {
	ColumnWithTable
	Direction OrderDirection
}

func NewColumnOrder(table TableType, column ColumnType, direction OrderDirection) ColumnOrder {
	return ColumnOrder{
		ColumnWithTable: ColumnWithTable{
			Table:  table,
			Column: column,
		},
		Direction: direction,
	}
}

func (c ColumnOrder) String() string {
	return fmt.Sprintf("%s %s", c.GetAlias(), c.Direction.Value)
}

func ORDER_BY(c ...ColumnOrder) Statement {
	columnsStr := ""
	for i, column := range c {
		columnsStr += StringWithWithoutComma(len(c), i, column.String())
	}

	stmt := fmt.Sprintf("ORDER BY %s", columnsStr)

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

// GROUP BY

func GROUP_BY(c ...ColumnWithTable) Statement {
	columnsStr := ""

	for i, column := range c {
		columnsStr += StringWithWithoutComma(len(c), i, column.GetAlias())
	}

	stmt := fmt.Sprintf("GROUP BY %s", columnsStr)

	return Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

// # LIMIT

func LIMIT(limit int) Statement {
	return Statement{
		SQL:  fmt.Sprintf("LIMIT %s", QUERY_ARG),
		Args: []interface{}{limit},
	}
}

// # OFFSET

func OFFSET(limit int) Statement {
	return Statement{
		SQL:  fmt.Sprintf("OFFSET %s", QUERY_ARG),
		Args: []interface{}{limit},
	}
}
