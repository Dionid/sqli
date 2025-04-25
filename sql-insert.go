package sqli

import "fmt"

func INSERT_INTO(t TableType, columns ...ColumnType) Statement {
	columnsStr := ""

	for i, column := range columns {
		if len(columns) > 1 {
			if i == len(columns)-1 {
				columnsStr += column.GetName()
			} else {
				columnsStr += column.GetName() + ", "
			}
		} else {
			columnsStr += column.GetName()
		}
	}

	if columnsStr == "" {
		return Statement{
			SQL:  fmt.Sprintf("INSERT INTO %s", t.GetName()),
			Args: []interface{}{},
		}
	} else {
		return Statement{
			SQL:  fmt.Sprintf("INSERT INTO %s (%s)", t.GetName(), columnsStr),
			Args: []interface{}{},
		}
	}
}

// TODO: Someday
//
// type SmartInsertValueSt[V any] struct {
// 	Column Column[V]
// 	Value V
// }

// func (s SmartInsertValueSt[V]) GetColumn() Column[V] {
// 	return s.Column
// }

// func (s SmartInsertValueSt[V]) GetValue() V {
// 	return s.Value
// }

// type SmartInsertValue[V any] interface {
// 	GetColumn() Column[V]
// 	GetValue() []V
// }

// func SMART_INSERT_VALUE[V any](c Column[V], value V) SmartInsertValueSt[V] {
// 	return SmartInsertValueSt[V]{
// 		Column: c,
// 		Value: value,
// 	}
// }

// type SmartInsertElementGroup[V any] []SmartInsertValue[V]

// func SMART_INSERT_INTO(t Table, elementGroupList ...SmartInsertElementGroup[any]) SqlWithArgs {
// 	columnsStr := ""
// 	valuesStr := ""
// 	args := []interface{}{}

// 	for i, elementGroup := range elementGroupList {
// 		columnsStr += "("
// 		for j, element := range elementGroup {
// 			switch e := element.(type) {
// 			case SmartInsertValue[any]:
// 				columnsStr += StringWithWithoutComma(len(elementGroup), j, e.GetColumn().GetName())
// 				valuesStr += StringWithWithoutComma(len(elementGroup), j, QUERY_ARG)
// 				args = append(args, e.GetValue())
// 				break;
// 			default:
// 				break;
// 			}
// 		}
// 		columnsStr += StringWithWithoutComma(len(elementGroupList), i, ")")
// 		valuesStr += StringWithWithoutComma(len(elementGroupList), i, ")")
// 	}

// 	return SqlWithArgs{
// 		SQL: fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", t.GetName(), columnsStr, valuesStr),
// 		Args: args,
// 	}
// }
