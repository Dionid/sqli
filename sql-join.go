package sqlification

import "fmt"

func LEFT_JOIN(table NameWithAliaser, on Statement) Statement {
	stmt := fmt.Sprintf("LEFT JOIN %s ON %s", table.GetNameWithAlias(), on.SQL)

	return Statement{
		SQL:  stmt,
		Args: on.Args,
	}
}

func RIGHT_JOIN(table NameWithAliaser, on Statement) Statement {
	stmt := fmt.Sprintf("RIGHT JOIN %s ON %s", table.GetNameWithAlias(), on.SQL)

	return Statement{
		SQL:  stmt,
		Args: on.Args,
	}
}

func INNER_JOIN(table NameWithAliaser, on Statement) Statement {
	stmt := fmt.Sprintf("INNER JOIN %s ON %s", table.GetNameWithAlias(), on.SQL)

	return Statement{
		SQL:  stmt,
		Args: on.Args,
	}
}

func FULL_JOIN(table NameWithAliaser, on Statement) Statement {
	stmt := fmt.Sprintf("FULL JOIN %s ON %s", table.GetNameWithAlias(), on.SQL)

	return Statement{
		SQL:  stmt,
		Args: on.Args,
	}
}

func CROSS_JOIN(table NameWithAliaser, on Statement) Statement {
	stmt := fmt.Sprintf("CROSS JOIN %s ON %s", table.GetNameWithAlias(), on.SQL)

	return Statement{
		SQL:  stmt,
		Args: on.Args,
	}
}
