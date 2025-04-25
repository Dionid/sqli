package sqli

import (
	"errors"
	"fmt"
	"strings"
)

type Column[V any] struct {
	ColumnName  string
	ColumnAlias string
	TableName   string
	TableAlias  string
}

func NewColumn[V any](table Table, columnName string) Column[V] {
	return Column[V]{
		ColumnName:  DoubleQuotes(columnName),
		ColumnAlias: DoubleQuotes(columnName),
		TableName:   DoubleQuotes(table.TableName),
		TableAlias:  DoubleQuotes(table.TableAlias),
	}
}

func NewColumnWithAlias[V any](table Table, columnName string, columnAlias string) Column[V] {
	return Column[V]{
		ColumnName:  DoubleQuotes(columnName),
		ColumnAlias: DoubleQuotes(columnAlias),
		TableName:   DoubleQuotes(table.TableName),
		TableAlias:  DoubleQuotes(table.TableAlias),
	}
}

func (c Column[V]) GetName() string {
	return c.ColumnName
}

func (c Column[V]) GetAlias() string {
	return c.ColumnAlias
}

func (c Column[V]) SetAlias(alias string) Column[V] {
	return Column[V]{
		ColumnName:  c.ColumnName,
		ColumnAlias: DoubleQuotes(alias),
		TableName:   c.TableName,
		TableAlias:  c.TableAlias,
	}
}

func (c Column[V]) GetNameAsAlias() string {
	if c.ColumnName == c.ColumnAlias {
		return DoubleQuotes(c.GetAlias())
	}
	return fmt.Sprintf(`%s AS %s`, c.GetName(), c.GetAlias())
}

func (c Column[V]) GetAliasWithTableAlias() string {
	if c.TableAlias == "" {
		return c.GetAlias()
	}
	return fmt.Sprintf(`%s.%s`, c.TableAlias, c.GetAlias())
}

func (c Column[V]) GetStatement() Statement {
	return Statement{
		SQL:  c.ColumnAlias,
		Args: []interface{}{},
	}
}

var AllColumn = Column[string]{"*", "*", "", ""}

type ColumnType interface {
	GetName() string
	GetAlias() string
	GetNameAsAlias() string
	GetAliasWithTableAlias() string
}

type Table struct {
	TableName  string
	TableAlias string
}

func NewTableSt(name string, alias string) Table {
	return Table{
		TableName:  fmt.Sprintf(`"%s"`, name),
		TableAlias: fmt.Sprintf(`"%s"`, alias),
	}
}

func (t Table) GetStatement() Statement {
	return Statement{
		SQL:  t.TableAlias,
		Args: []interface{}{},
	}
}

func (t Table) GetAliasAsStatement() Statement {
	return Statement{
		SQL:  t.TableAlias,
		Args: []interface{}{},
	}
}

func (t Table) GetName() string {
	return DoubleQuotes(t.TableName)
}

func (t Table) As(alias string) TableType {
	t.TableAlias = DoubleQuotes(alias)

	return t
}

func (t Table) GetNameWithAlias() string {
	return fmt.Sprintf(`%s AS %s`, t.TableName, t.TableAlias)
}

func (t Table) GetAlias() string {
	return t.TableAlias
}

func (t Table) AllColumns() Column[string] {
	return AllColumn
}

type TableType interface {
	GetAlias() string
	GetName() string
	GetNameWithAlias() string
}

type ColumnWithTable struct {
	Table  TableType
	Column ColumnType
}

func NewColumnWithTable(table TableType, column ColumnType) ColumnWithTable {
	// TODO: # Check that column exists
	// ...
	return ColumnWithTable{
		Table:  table,
		Column: column,
	}
}

func (c ColumnWithTable) GetAlias() string {
	return fmt.Sprintf("%s.%s", c.Table.GetAlias(), c.Column.GetAlias())
}

// # STATEMENT

type Statement struct {
	SQL  string
	Args []interface{}
}

func NewStatement(sql string, args ...interface{}) Statement {
	return Statement{
		SQL:  sql,
		Args: args,
	}
}

func NewEmptyStatement() Statement {
	return Statement{
		SQL:  "",
		Args: []interface{}{},
	}
}

func (s Statement) String() string {
	return s.SQL
}

func (s Statement) GetSQL() string {
	return s.SQL
}

func (s Statement) GetArguments() []interface{} {
	return s.Args
}

func (s Statement) GetStatement() Statement {
	return s
}

func (s Statement) GetName() string {
	return s.SQL
}

func (s Statement) GetAlias() string {
	return s.SQL
}

func (s Statement) GetNameAsAlias() string {
	return s.SQL
}

func (s Statement) GetAliasWithTableAlias() string {
	return s.SQL
}

type StatementType interface {
	GetStatement() Statement
}

func StatementFromTableAlias(table TableType) Statement {
	return Statement{
		SQL:  table.GetAlias(),
		Args: []interface{}{},
	}
}

func AggregateStatements(queries ...Statement) Statement {
	sqlQuery := ""
	args := []interface{}{}

	for i, query := range queries {
		if i == 0 {
			sqlQuery += query.SQL
		} else {
			sqlQuery += (" " + query.SQL)
		}
		args = append(args, query.Args...)
	}

	return Statement{
		SQL:  sqlQuery,
		Args: args,
	}
}

func SetArgsSequence(val Statement) Statement {
	sqlWithArgs := ""
	argCounter := 1

	for i, parts := range strings.Split(val.SQL, QUERY_ARG) {
		if i == 0 {
			sqlWithArgs += parts
		} else {
			sqlWithArgs += fmt.Sprintf("$%d", argCounter)
			argCounter++
			sqlWithArgs += parts
		}
	}

	return Statement{
		SQL:  sqlWithArgs,
		Args: val.Args,
	}
}

// # QUERY

type QueryOptions struct {
	CheckForDeleteWhere bool
	CheckForUpdateWhere bool
	EndWithSemicolon    bool
}

func QueryWithOptions(options QueryOptions, exprs ...Statement) (Statement, error) {
	queryAggregate := SetArgsSequence(AggregateStatements(exprs...))

	if options.EndWithSemicolon {
		queryAggregate.SQL += ";"
	}

	checkUpdate := options.CheckForUpdateWhere && strings.Contains(queryAggregate.SQL, `UPDATE `) && strings.Contains(queryAggregate.SQL, ` SET `)
	checkDelete := options.CheckForDeleteWhere && strings.Contains(queryAggregate.SQL, "DELETE FROM")

	if checkUpdate || checkDelete {
		if !strings.Contains(queryAggregate.SQL, "WHERE") {
			return queryAggregate, errors.New("WHERE clause is required for UPDATE and DELETE queries")
		}
	}

	return queryAggregate, nil
}

func Query(exprs ...Statement) (Statement, error) {
	return QueryWithOptions(
		QueryOptions{
			CheckForDeleteWhere: true,
			CheckForUpdateWhere: true,
			EndWithSemicolon:    true,
		},
		exprs...,
	)
}

func QueryMust(exprs ...Statement) Statement {
	stmt, err := QueryWithOptions(
		QueryOptions{
			CheckForDeleteWhere: true,
			CheckForUpdateWhere: true,
			EndWithSemicolon:    true,
		},
		exprs...,
	)

	if err != nil {
		panic(err)
	}

	return stmt
}

func SubQuery(exprs ...Statement) Statement {
	sp := AggregateStatements(exprs...)
	sp.SQL = fmt.Sprintf("(%s)", sp.SQL)
	return sp
}

type NameWithAliaser interface {
	GetNameWithAlias() string
}

type CompareOperator struct {
	Value string
}

func (c CompareOperator) String() string {
	return c.Value
}

func Compare[V any](c Column[V], operator CompareOperator, value V) Statement {
	return Statement{
		SQL:  fmt.Sprintf("%s.%s %s %s", c.TableAlias, c.GetName(), operator, QUERY_ARG),
		Args: []interface{}{value},
	}
}

func CompareColumns[A any, B any](a Column[A], operator CompareOperator, b Column[B]) Statement {
	return Statement{
		SQL:  fmt.Sprintf("%s.%s %s %s.%s", a.TableAlias, a.GetName(), operator, b.TableAlias, b.GetName()),
		Args: []interface{}{},
	}
}
