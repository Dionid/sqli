package sqlification

import (
	"fmt"
)

func WHERE(c Statement) Statement {
	return Statement{
		SQL:  fmt.Sprintf("WHERE %s", c),
		Args: c.Args,
	}
}

func AND(c ...Statement) Statement {
	sql := "("
	args := []interface{}{}

	for i, v := range c {
		if i == 0 {
			sql += v.SQL
		} else {
			sql += fmt.Sprintf(" AND %s", v.SQL)
		}
		args = append(args, v.Args...)
	}

	sql += ")"

	return Statement{
		SQL:  sql,
		Args: args,
	}
}

func OR(c ...Statement) Statement {
	sql := "("
	args := []interface{}{}

	for i, v := range c {
		if i == 0 {
			sql += v.SQL
		} else {
			sql += fmt.Sprintf(" OR %s", v.SQL)
		}
		args = append(args, v.Args...)
	}

	sql += ")"

	return Statement{
		SQL:  sql,
		Args: args,
	}
}

var Equal = CompareOperator{
	Value: "=",
}

func EQUAL[V any](c Column[V], value V) Statement {
	return Compare(c, Equal, value)
}

func EQUAL_COLUMNS[A any, B any](a Column[A], b Column[B]) Statement {
	return CompareColumns[A, B](a, Equal, b)
}

var NotEqual = CompareOperator{
	Value: "!=",
}

func NOT_EQUAL[V any](c Column[V], value V) Statement {
	return Compare(c, NotEqual, value)
}

var In = CompareOperator{
	Value: "IN",
}

func IN[V any](c Column[V], value V) Statement {
	return Compare(c, In, value)
}

var NotIn = CompareOperator{
	Value: "NOT IN",
}

func NOT_IN[V any](c Column[V], value V) Statement {
	return Compare(c, NotIn, value)
}

var Less = CompareOperator{
	Value: "<",
}

func LESS[V any](c Column[V], value V) Statement {
	return Compare(c, Less, value)
}

var Greater = CompareOperator{
	Value: ">",
}

func GREATER[V any](c Column[V], value V) Statement {
	return Compare(c, Greater, value)
}

var LessOrEqual = CompareOperator{
	Value: "<=",
}

func LESS_OR_EQUAL[V any](c Column[V], value V) Statement {
	return Compare(c, LessOrEqual, value)
}

var GreaterOrEqual = CompareOperator{
	Value: ">=",
}

func GREATER_OR_EQUAL[V any](c Column[V], value V) Statement {
	return Compare(c, GreaterOrEqual, value)
}

var Like = CompareOperator{
	Value: "LIKE",
}

func LIKE[V any](c Column[V], value V) Statement {
	return Compare(c, Like, value)
}

var ILike = CompareOperator{
	Value: "ILIKE",
}

func ILIKE[V any](c Column[V], value V) Statement {
	return Compare(c, ILike, value)
}
