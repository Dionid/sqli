![logo.png](./logo.png)

# SQLi

Type-safe generative SQL Query Builder based on you DB schema.

```go

// SQLi generates Predefined Queries

newUserId := uuid.MustParse("ebb5dd71-1214-40dc-b602-bb9af74b3aae")

id, err := InsertIntoUserReturningID( // This is a generated function
    ctx,
    db,
    yourdb.InsertableUserModel{ 
        ID: newUserId,
        Name: "User 1",
    },
)

println(id) // ebb5dd71-1214-40dc-b602-bb9af74b3aae

// ... and you can also use Dynamic Queries

insertQuery, err := sqli.Query(
    sqli.INSERT_INTO(yourdb.User),
    sqli.VALUES(
        sqli.ValueSet(
            VALUE(yourdb.User.ID, newUserId), // VALUE function will validate the value UUID type
            VALUE(yourdb.User.Name, "User 1"), // And this will validate the value string type
        ),
    ),
    sqli.RETURNING(
        yourdb.User.ID,
    ),
)

println(insertQuery.SQL) // `INSERT INTO "user" (VALUES ($1, $2));`
println(insertQuery.Args) // [ebb5dd71-1214-40dc-b602-bb9af74b3aae, User 1]

row := db.QueryRowxContext(ctx, insertQuery.SQL, insertQuery.Args...)

var id uuid.UUID
err = row.Scan(&id)

println(id) // ebb5dd71-1214-40dc-b602-bb9af74b3aae

```

# Install

## By go install

```shell
go install github.com/Dionid/sqli/cmd/sqli@latest
```

## Releases

Download from [Release](https://github.com/Dionid/sqli/releases)

## From sources

```shell
git clone git@github.com:Dionid/sqli.git
cd sqli
make build
# you will find executables in ./dist folder
```

# Features

SQLi generates:

1. Constants
1. Predefined Queries
3. Dynamic Queries

## Constants

All tables as types and names can be found in `constants.sqli.go`

```go
type TablesSt struct {
	Office         string `json:"office" db:"office"`
	OfficeUser     string `json:"office_user" db:"office_user"`
	User           string `json:"user" db:"user"`
}

var Tables = TablesSt{
	Office:         "office",
	OfficeUser:     "office_user",
	User:           "user",
}

// Named "T" for shortness
var T = Tables
```

## Predefined Queries

This functions are generated for each table in your DB schema and capture
most common queries, like Insert, Update, Delete, Select by Primary key / Unique key / etc.

Lets look at the example of generated function for `user` table:

```go

newUserId := uuid.MustParse("ebb5dd71-1214-40dc-b602-bb9af74b3aae")

id, err := InsertIntoUserReturningID(
	ctx,
	db,
    InsertableUserModel{
        ID: newUserId,
        Name: "User 1",
    },
)

println(id) // ebb5dd71-1214-40dc-b602-bb9af74b3aae

// Not lets select it by primary key

userByPrimaryKey, err := SelectUserByID(
    ctx,
    db,
    id,
)

println(userByPrimaryKey) // {ID: "ebb5dd71-1214-40dc-b602-bb9af74b3aae", Name: "User 1"}

// Now lets update it

err = UpdateUserByID(
    ctx,
    db,
    id,
    UpdatableOfficeModel{
        Name: "Updated User 1",
    }
)

// And delete

err = DeleteFromUserByID(
    ctx,
    db,
    id,
)

```

### Insert

1. Returning Result
1. Returning All
1. Returning Primary key
1. Returning Unique key

### Select By

1. Primary
1. By Primary compound
1. Sequence

### Delete By

1. Primary
1. By Primary compound
1. Sequence

### Update

1. By Primary
1. By Primary compound
1. By Sequence

# Dynamic Type-safe Queries

This is a dynamic query builder, that allows you to build queries
in a type-safe way, using the generated constants and functions.

```go

import (
    .   "github.com/Dionid/sqli"
    .   "github.com/Dionid/sqli/examples/pgdb/db"
)

// Insert user

newUserId := uuid.MustParse("ebb5dd71-1214-40dc-b602-bb9af74b3aae")

insertQuery, err := Query(
    INSERT_INTO(User),
    VALUES(
        ValueSet(
            VALUE(User.ID, newUserId), // VALUE function will validate the value UUID type
            VALUE(User.Name, "User 1"), // And this will validate the value string type
        ),
    ),
    RETURNING(
        User.ID,
    ),
)

// `insertQuery` will have raw SQL and raw arguments, that can be used to execute the query

println(insertQuery.SQL) // `INSERT INTO "user" (VALUES ($1, $2));`
println(insertQuery.Args) // [ebb5dd71-1214-40dc-b602-bb9af74b3aae, User 1]

// Now lets execute the query
row := db.QueryRowxContext(ctx, insertQuery.SQL, insertQuery.Args...)

var id uuid.UUID
err = row.Scan(&id)

println(id) // ebb5dd71-1214-40dc-b602-bb9af74b3aae

// Select it by primary key

selectQuery, err := Query(
    SELECT(
        User.AllColumns(), // *
    ),
    FROM(User),
    WHERE(
        EQUAL(User.ID, id),
    ),
)

println(selectQuery.SQL) // SELECT * FROM "user" AS "user" WHERE "user"."id" = $1;
println(selectQuery.Args) // [ebb5dd71-1214-40dc-b602-bb9af74b3aae]

row := db.QueryRowxContext(ctx, selectQuery.SQL, selectQuery.Args...)
user := &UserModel{} // Also generated by SQLi
err = row.Scan(
    &user.ID,
    &user.Name,
)
println(user) // {ID: "ebb5dd71-1214-40dc-b602-bb9af74b3aae", Name: "User 1"}

// Now lets update it

query, err := Query(
    UPDATE(User),
    SET(
        SET_VALUE(User.Name, "Updated User 1"),
    ),
    WHERE(
        EQUAL(User.ID, id),
    ),
)

println(query.SQL) // UPDATE "user" SET "user"."name" = $1 WHERE "user"."id" = $2;
println(query.Args) // [Updated User 1, ebb5dd71-1214-40dc-b602-bb9af74b3aae]

row := db.ExecContext(ctx, query.SQL, query.Args...)
println(row.RowsAffected()) // 1

// And delete

query, err := Query(
    DELETE_FROM(User),
    WHERE(
        EQUAL(User.ID, id),
    ),
)

println(query.SQL) // DELETE FROM "user" WHERE "user"."id" = $1;
println(query.Args) // [ebb5dd71-1214-40dc-b602-bb9af74b3aae]

row := db.ExecContext(ctx, query.SQL, query.Args...)
println(row.RowsAffected()) // 1

```

# Examples

For more examples, see the [examples](./examples) folder.

# What is the difference between SQLi and other query builders?

## Generative

## Type-safe

Every function is generated for each table in your DB schema and typed according to the table schema.

```go

EQUAL(User.ID, 123) // This will not compile, because ID is UUID
EQUAL(User.ID, uuid.MustParse("ebb5dd71-1214-40dc-b602-bb9af74b3aae")) // This will compile

```

## Extensible

Most Query Builders uses dot notation to build queries, like `db.Table("user").Where("id = ?", id)`,
but SQLi uses functional approach `Query(SELECT(User.ID), FROM(User) WHERE(EQUAL(User.ID, id))` where
each function needs to return a `Statement` struct, that contains SQL and arguments.

That gives us ability to extend the library and add new functions, like `JSON_AGG`, `SUM`, `COUNT`, etc.
WITHOUT even commiting to the library itself.

Example:

We got some database that has operators like `MERGE table WHERE ... COLLISION free | restricted`, but we
don't has this functions in SQLi, so we can create our own functions and use them in the query builder:

```go
func MERGE(table NameWithAliaser) string {
    stmt := fmt.Sprintf("FROM %s", table.GetNameWithAlias())

    return sqli.Statement{
		SQL:  stmt,
		Args: []interface{}{},
	}
}

func COLLISION_FREE() string {
    return sqli.Statement{
		SQL:  "COLLISION free",
		Args: []interface{}{},
	}
}

func COLLISION_RESTRICTED() string {
    return sqli.Statement{
		SQL:  "COLLISION restricted",
		Args: []interface{}{},
	}
}

func main() {
    query := Query(
        MERGE("table"),
        WHERE(
            EQUAL("id", 1),
        ),
        COLLISION_FREE(),
    )
    fmt.Println(query.SQL) // MERGE table WHERE id = 1 COLLISION free
}
```

So you don't even need to wait for the library to implement this functions, you can do it yourself.

# TODO

1. Upsert
1. SUM
1. JSON_AGG
1. InsertOnConflict
1. CopySimple
1. Add pgx types
1. Safe-mode (validating every field)