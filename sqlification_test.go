package sqlification_test

import (
	"testing"

	. "github.com/Dionid/sqlification"
	. "github.com/Dionid/sqlification/dbs/crmdb"
	"github.com/google/uuid"
)

func TestSimpleUpdate(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE(UserTable.Password, "test"),
		),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = $1 WHERE "user"."id" = $2;` {
		t.Errorf("Query is not correct: %s", query.SQL)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != id {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestSimpleUpdateError(t *testing.T) {
	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE(UserTable.Password, "test"),
		),
	)
	if err == nil {
		t.Errorf("Error is nil: %s", query)
	}
}

func TestUpdateMoreArgs(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE(UserTable.Password, "test"),
			SET_VALUE(UserTable.Email, "some@mail.com"),
		),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = $1, "email" = $2 WHERE "user"."id" = $3;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 3 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != `some@mail.com` {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}

	if query.Args[2] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}
}

func TestUpdateReturningAll(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE(UserTable.Password, "test"),
		),
		WHERE(EQUAL(UserTable.ID, id)),
		RETURNING(UserTable.AllColumns()),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = $1 WHERE "user"."id" = $2 RETURNING *;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != id {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestUpdateReturning(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE(UserTable.Password, "test"),
		),
		WHERE(EQUAL(UserTable.ID, id)),
		RETURNING(UserTable.ID),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = $1 WHERE "user"."id" = $2 RETURNING "id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != id {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestUpdateWith(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	tmpTableBase := NewTableSt("tmp", "tmp")

	tmpTable := struct {
		Table

		Password Column[string]
		Email    Column[string]
	}{
		tmpTableBase,
		NewColumn[string](tmpTableBase, "password"),
		NewColumn[string](tmpTableBase, "email"),
	}

	query, err := Query(
		WITH(
			AS(
				StatementFromTableAlias(tmpTable),
				SubQuery(
					SELECT(
						UserTable.Password,
						UserTable.Email,
					),
					FROM(UserTable),
					WHERE(EQUAL(UserTable.ID, id)),
				),
			),
		),
		UPDATE(UserTable),
		SET(
			SET_VALUE_COLUMN(UserTable.Password, tmpTable.Password),
			SET_VALUE_COLUMN(UserTable.Email, tmpTable.Email),
		),
		FROM(tmpTable),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `WITH "tmp" AS (SELECT "user"."password", "user"."email" FROM "user" AS "user" WHERE "user"."id" = $1) UPDATE "user" SET "password" = "tmp"."password", "email" = "tmp"."email" FROM "tmp" AS "tmp" WHERE "user"."id" = $2;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}

	if query.Args[1] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}
}

func TestUpdateFromValues(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	tmpTableBase := NewTableSt("tmp", "tmp")

	tmpTable := struct {
		Table

		Password Column[string]
		Email    Column[string]
	}{
		tmpTableBase,
		NewColumn[string](tmpTableBase, "password"),
		NewColumn[string](tmpTableBase, "email"),
	}

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE_COLUMN(UserTable.Password, tmpTable.Password),
			SET_VALUE_COLUMN(UserTable.Email, tmpTable.Email),
		),
		FROM_RAW(
			AS(
				VALUES(
					ValuesSetSt{
						VALUE(tmpTable.Password, "test"),
						VALUE(tmpTable.Email, "mail@test.com"),
					},
				),
				TABLE_WITH_COLUMNS(
					tmpTable,
					tmpTable.Password,
					tmpTable.Email,
				),
			),
		),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = "tmp"."password", "email" = "tmp"."email" FROM (VALUES ($1, $2)) AS "tmp"("password", "email") WHERE "user"."id" = $3;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 3 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != `mail@test.com` {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}

	if query.Args[2] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}
}

func TestUpdateFromOtherTable(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	tmpTable := UserTable.As("usr")

	query, err := Query(
		UPDATE(UserTable),
		SET(
			SET_VALUE_COLUMN(UserTable.Password, tmpTable.Password),
			SET_VALUE_COLUMN(UserTable.Email, tmpTable.Email),
		),
		FROM_RAW(
			AS(
				SubQuery(
					SELECT(
						tmpTable.Password,
						tmpTable.Email,
					),
					FROM(tmpTable),
					WHERE(
						EQUAL(tmpTable.ID, id),
					),
				),
				tmpTable,
			),
		),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `UPDATE "user" SET "password" = "usr"."password", "email" = "usr"."email" FROM (SELECT "usr"."password", "usr"."email" FROM "user" AS "usr" WHERE "usr"."id" = $1) AS "usr" WHERE "user"."id" = $2;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != id {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestUpdateFromWith(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	tmpTableBase := NewTableSt("tmp", "tmp")

	tmpTable := struct {
		Table

		Password Column[string]
		Email    Column[string]
	}{
		tmpTableBase,
		NewColumn[string](tmpTableBase, "password"),
		NewColumn[string](tmpTableBase, "email"),
	}

	query, err := Query(
		WITH(
			AS(
				TABLE_WITH_COLUMNS(
					tmpTable,
					tmpTable.Password,
					tmpTable.Email,
				),
				VALUES(
					ValuesSetSt{
						VALUE(tmpTable.Password, "test"),
						VALUE(tmpTable.Email, "mail@test.com"),
					},
				),
			),
		),
		UPDATE(UserTable),
		SET(
			SET_VALUE_COLUMN(UserTable.Password, tmpTable.Password),
			SET_VALUE_COLUMN(UserTable.Email, tmpTable.Email),
		),
		FROM(tmpTable),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `WITH "tmp"("password", "email") AS (VALUES ($1, $2)) UPDATE "user" SET "password" = "tmp"."password", "email" = "tmp"."email" FROM "tmp" AS "tmp" WHERE "user"."id" = $3;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 3 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != `mail@test.com` {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}

	if query.Args[2] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}
}

func TestSimpleDelete(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		DELETE_FROM(UserTable),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `DELETE FROM "user" WHERE "user"."id" = $1;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 1 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}
}

func TestSimpleDeleteError(t *testing.T) {
	query, err := Query(
		DELETE_FROM(UserTable),
	)
	if err == nil {
		t.Errorf("Error is nil: %s", query)
	}
}

func TestInsertSimple(t *testing.T) {
	query, err := Query(
		INSERT_INTO(UserTable),
		VALUES(
			ValueSet(
				VALUE(UserTable.Password, "test"),
				VALUE(UserTable.Email, "mail@test.com"),
			),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `INSERT INTO "user" (VALUES ($1, $2));` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != `mail@test.com` {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestInsertWithColumns(t *testing.T) {
	query, err := Query(
		INSERT_INTO(UserTable, UserTable.Password, UserTable.Email),
		VALUES(
			ValuesSetSt{
				VALUE(UserTable.Password, "test"),
				VALUE(UserTable.Email, "mail@test.com"),
			},
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `INSERT INTO "user" ("password", "email") (VALUES ($1, $2));` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != `test` {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != `mail@test.com` {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestSimpleSelect(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Password,
			UserTable.Email,
		),
		FROM(UserTable),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."password", "user"."email" FROM "user" AS "user";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectAll(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectWhere(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(EQUAL(UserTable.ID, id)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE "user"."id" = $1;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 1 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}
}

func TestSelectWhereOr(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(
			OR(
				EQUAL(UserTable.ID, id),
				EQUAL(UserTable.Email, "mail@test.com"),
			),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE ("user"."id" = $1 OR "user"."email" = $2);` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != "mail@test.com" {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestSelectWhereAnd(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(
			AND(
				EQUAL(UserTable.ID, id),
				EQUAL(UserTable.Email, "mail@test.com"),
			),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE ("user"."id" = $1 AND "user"."email" = $2);` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 2 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != "mail@test.com" {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}
}

func TestSelectWhereOrAnd(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(
			OR(
				AND(
					EQUAL(UserTable.ID, id),
					EQUAL(UserTable.Email, "mail@test.com"),
				),
				AND(
					EQUAL(UserTable.ID, id),
					EQUAL(UserTable.Email, "mail2@test.com"),
				),
			),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE (("user"."id" = $1 AND "user"."email" = $2) OR ("user"."id" = $3 AND "user"."email" = $4));` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 4 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}

	if query.Args[1] != "mail@test.com" {
		t.Errorf("Arg is not correct: %s", query.Args[1])
	}

	if query.Args[2] != id {
		t.Errorf("Arg is not correct: %s", query.Args[2])
	}

	if query.Args[3] != "mail2@test.com" {
		t.Errorf("Arg is not correct: %s", query.Args[3])
	}
}

func TestSelectOrderBy(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(EQUAL(UserTable.ID, id)),
		ORDER_BY(NewColumnOrder(UserTable, UserTable.CreatedAt, DESC)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE "user"."id" = $1 ORDER BY "user"."created_at" DESC;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 1 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}
}

func TestSelectGroupBy(t *testing.T) {
	id := uuid.MustParse("c153f078-9ca1-455e-90ff-dc9975948259")

	query, err := Query(
		SELECT(
			UserTable.AllColumns(),
		),
		FROM(UserTable),
		WHERE(EQUAL(UserTable.ID, id)),
		GROUP_BY(NewColumnWithTable(UserTable, UserTable.Email)),
		ORDER_BY(NewColumnOrder(UserTable, UserTable.CreatedAt, DESC)),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT * FROM "user" AS "user" WHERE "user"."id" = $1 GROUP BY "user"."email" ORDER BY "user"."created_at" DESC;` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 1 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}

	if query.Args[0] != id {
		t.Errorf("Arg is not correct: %s", query.Args[0])
	}
}

func TestSelectLeftJoin(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Name,
		),
		FROM(
			UserTable,
		),
		LEFT_JOIN(
			OfficeUserTable,
			EQUAL_COLUMNS(OfficeUserTable.UserID, UserTable.ID),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."name" FROM "user" AS "user" LEFT JOIN "office_user" AS "office_user" ON "office_user"."user_id" = "user"."id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectRightJoin(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Name,
		),
		FROM(
			UserTable,
		),
		RIGHT_JOIN(
			OfficeUserTable,
			EQUAL_COLUMNS(OfficeUserTable.UserID, UserTable.ID),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."name" FROM "user" AS "user" RIGHT JOIN "office_user" AS "office_user" ON "office_user"."user_id" = "user"."id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectInnerJoin(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Name,
		),
		FROM(
			UserTable,
		),
		INNER_JOIN(
			OfficeUserTable,
			EQUAL_COLUMNS(OfficeUserTable.UserID, UserTable.ID),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."name" FROM "user" AS "user" INNER JOIN "office_user" AS "office_user" ON "office_user"."user_id" = "user"."id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectFullJoin(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Name,
		),
		FROM(
			UserTable,
		),
		FULL_JOIN(
			OfficeUserTable,
			EQUAL_COLUMNS(OfficeUserTable.UserID, UserTable.ID),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."name" FROM "user" AS "user" FULL JOIN "office_user" AS "office_user" ON "office_user"."user_id" = "user"."id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectCrossJoin(t *testing.T) {
	query, err := Query(
		SELECT(
			UserTable.Name,
		),
		FROM(
			UserTable,
		),
		CROSS_JOIN(
			OfficeUserTable,
			EQUAL_COLUMNS(OfficeUserTable.UserID, UserTable.ID),
		),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT "user"."name" FROM "user" AS "user" CROSS JOIN "office_user" AS "office_user" ON "office_user"."user_id" = "user"."id";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

func TestSelectCount(t *testing.T) {
	query, err := Query(
		SELECT(
			COUNT(UserTable.ID),
		),
		FROM(UserTable),
	)
	if err != nil {
		t.Errorf("Error is not nil: %s", err)
	}

	if query.SQL != `SELECT COUNT("user"."id") FROM "user" AS "user";` {
		t.Errorf("Query is not correct: %s", query)
	}

	if len(query.Args) != 0 {
		t.Errorf("Args is not correct: %d", len(query.Args))
	}
}

// TODO: # SELECT SUM
// ...

// TODO: # SELECT FROM SUB QUERY
// ...

// TODO: # SELECT JSON_AGG
// ...

// TODO: # InsertReturning
// ...

// TODO: # InsertOnConflict
// ...

// TODO: # CopySimple
// ...
