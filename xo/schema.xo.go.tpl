{{ define "enum" }}
{{- $e := .Data -}}
// {{ $e.GoName }} is the '{{ $e.SQLName }}' enum type from schema '{{ schema }}'.
type {{ $e.GoName }} uint16

// {{ $e.GoName }} values.
const (
{{ range $e.Values -}}
	// {{ $e.GoName }}{{ .GoName }} is the '{{ .SQLName }}' {{ $e.SQLName }}.
	{{ $e.GoName }}{{ .GoName }} {{ $e.GoName }} = {{ .ConstValue }}
{{ end -}}
)

// String satisfies the [fmt.Stringer] interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) String() string {
	switch {{ short $e.GoName }} {
{{ range $e.Values -}}
	case {{ $e.GoName }}{{ .GoName }}:
		return "{{ .SQLName }}"
{{ end -}}
	}
	return fmt.Sprintf("{{ $e.GoName }}(%d)", {{ short $e.GoName }})
}

// MarshalText marshals [{{ $e.GoName }}] into text.
func ({{ short $e.GoName }} {{ $e.GoName }}) MarshalText() ([]byte, error) {
	return []byte({{ short $e.GoName }}.String()), nil
}

// UnmarshalText unmarshals [{{ $e.GoName }}] from text.
func ({{ short $e.GoName }} *{{ $e.GoName }}) UnmarshalText(buf []byte) error {
	switch str := string(buf); str {
{{ range $e.Values -}}
	case "{{ .SQLName }}":
		*{{ short $e.GoName }} = {{ $e.GoName }}{{ .GoName }}
{{ end -}}
	default:
		return ErrInvalid{{ $e.GoName }}(str)
	}
	return nil
}

// Value satisfies the [driver.Valuer] interface.
func ({{ short $e.GoName }} {{ $e.GoName }}) Value() (driver.Value, error) {
	return {{ short $e.GoName }}.String(), nil
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ short $e.GoName }} *{{ $e.GoName }}) Scan(v interface{}) error {
	switch x := v.(type) {
	case []byte:
		return {{ short $e.GoName }}.UnmarshalText(x)
	case string:
		return {{ short $e.GoName }}.UnmarshalText([]byte(x))
	}
	return ErrInvalid{{ $e.GoName }}(fmt.Sprintf("%T", v))
}

{{ $nullName := (printf "%s%s" "Null" $e.GoName) -}}
{{- $nullShort := (short $nullName) -}}
// {{ $nullName }} represents a null '{{ $e.SQLName }}' enum for schema '{{ schema }}'.
type {{ $nullName }} struct {
	{{ $e.GoName }} {{ $e.GoName }}
	// Valid is true if [{{ $e.GoName }}] is not null.
	Valid bool
}

// Value satisfies the [driver.Valuer] interface.
func ({{ $nullShort }} {{ $nullName }}) Value() (driver.Value, error) {
	if !{{ $nullShort }}.Valid {
		return nil, nil
	}
	return {{ $nullShort }}.{{ $e.GoName }}.Value()
}

// Scan satisfies the [sql.Scanner] interface.
func ({{ $nullShort }} *{{ $nullName }}) Scan(v interface{}) error {
	if v == nil {
		{{ $nullShort }}.{{ $e.GoName }}, {{ $nullShort }}.Valid = 0, false
		return nil
	}
	err := {{ $nullShort }}.{{ $e.GoName }}.Scan(v)
	{{ $nullShort }}.Valid = err == nil
	return err
}

// ErrInvalid{{ $e.GoName }} is the invalid [{{ $e.GoName }}] error.
type ErrInvalid{{ $e.GoName }} string

// Error satisfies the error interface.
func (err ErrInvalid{{ $e.GoName }}) Error() string {
	return fmt.Sprintf("invalid {{ $e.GoName }}(%s)", string(err))
}
{{ end }}

{{ define "foreignkey" }}
{{- $i := .Data -}}


{{end}}

{{ define "procs" }}
{{- $ps := .Data -}}
{{- range $p := $ps -}}
// {{ func_name_context $p }} calls the stored {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func_context $p }} {
{{- if and (driver "mysql") (eq $p.Type "procedure") (not $p.Void) }}
	// At the moment, the Go MySQL driver does not support stored procedures
	// with out parameters
	return {{ zero $p.Returns }}, fmt.Errorf("unsupported")
{{- else }}
	// call {{ schema $p.SQLName }}
	{{ sqlstr "proc" $p }}
	// run
{{- if not $p.Void }}
{{- range $p.Returns }}
	var {{ check_name .GoName }} {{ type .Type }}
{{- end }}
	logf(sqlstr, {{ params $p.Params false }})
{{- if and (driver "sqlserver" "oracle") (eq $p.Type "procedure")}}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if err := {{ db "QueryRowx" $p }}.Scan({{ names "&" $p.Returns }}); err != nil {
{{- end }}
		return {{ zero $p.Returns }}, logerror(err)
	}
	return {{ range $p.Returns }}{{ check_name .GoName }}, {{ end }}nil
{{- else }}
	logf(sqlstr)
{{- if driver "sqlserver" "oracle" }}
	if _, err := {{ db_named "Exec" $p }}; err != nil {
{{- else }}
	if _, err := {{ db "Exec" $p }}; err != nil {
{{- end }}
		return logerror(err)
	}
	return nil
{{- end }}
{{- end }}
}

{{ if context_both -}}
// {{ func_name $p }} calls the {{ $p.Type }} '{{ $p.Signature }}' on db.
{{ func $p }} {
	return {{ func_name_context $p }}({{ names_all "" "context.Background()" "db" $p.Params }})
}
{{- end -}}
{{- end }}
{{ end }}

{{ define "typedef" }}
{{- $t := .Data -}}

type {{ $t.GoName }}TableSt struct {
	sqli.Table
	{{ range $t.Fields -}}
		{{ .GoName }} sqli.Column[{{ .Type }}]
	{{ end }}
}

func (t {{ $t.GoName }}TableSt) As(alias string) {{ $t.GoName }}TableSt {
	t.Table.TableAlias = fmt.Sprintf(`"%s"`, alias)
	{{ range $t.Fields -}}
		t.{{ .GoName }} = sqli.NewColumnWithAlias[{{ .Type }}](t.Table, t.{{ .GoName }}.ColumnName, t.{{ .GoName }}.ColumnAlias)
	{{ end }}

	return t
}

var {{ $t.GoName }}TableBase = sqli.Table{
	TableName: `"{{ $t.SQLName }}"`,
	TableAlias: `"{{ $t.SQLName }}"`,
}

var {{ $t.GoName }}Table = {{ $t.GoName }}TableSt{
	Table: {{ $t.GoName }}TableBase,
{{ range $t.Fields -}}
	{{ .GoName }}: sqli.NewColumn[{{ .Type }}]({{ $t.GoName }}TableBase, `"{{ .SQLName }}"`),
{{ end }}
}

// # Constants

{{ range $t.Fields -}}
type {{ $t.GoName }}{{ .GoName }}CT = {{ .Type }}
{{ end }}

{{ range $t.Fields -}}
const {{ $t.GoName }}{{ .GoName }}CN = `"{{ .SQLName }}"`
{{ end }}

// # Model

type {{ $t.GoName }}Model struct {
	{{ range $t.Fields -}}
		{{ .GoName }} {{ .Type }} `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
	{{ end }}
}

func New{{ $t.GoName }}Model(
{{ range $t.Fields -}}
	{{ .GoName }} {{ .Type }},
{{ end }}
) *{{ $t.GoName }}Model {
	return &{{ $t.GoName }}Model{
		{{ range $t.Fields -}}
			{{ .GoName }}: {{ .GoName }},
		{{ end }}
	}
}

// ## Insertable

type Insertable{{ $t.GoName }}Model struct {
	{{ range $t.Fields -}}
		{{- if not .IsSequence -}}
			{{- if eq .Default "" -}}
				{{ .GoName }} {{ .Type }} `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
			{{ else -}}
				{{ .GoName }} {{ .Type }} `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
			{{ end -}}
		{{- end -}}
	{{- end -}}
}

func NewInsertable{{ $t.GoName }}Model(
{{ range $t.Fields -}}
	{{- if not .IsSequence -}}
		{{- if eq .Default "" -}}
			{{ .GoName }} {{ .Type }},
		{{ else -}}
			{{ .GoName }} {{ .Type }},
		{{ end -}}
	{{- end -}}
{{- end -}}
) *Insertable{{ $t.GoName }}Model {
	return &Insertable{{ $t.GoName }}Model{
		{{ range $t.Fields -}}
			{{- if not .IsSequence -}}
				{{- if eq .Default "" -}}
					{{ .GoName }}: {{ .GoName }},
				{{ else -}}
					{{ .GoName }}: {{ .GoName }},
				{{ end -}}
			{{- end -}}
		{{- end -}}
	}
}

func InsertInto{{ $t.GoName }}Table(
	ctx context.Context,
	db DB,
	modelsList ...*Insertable{{ $t.GoName }}Model,
) (sql.Result, error) {
	if modelsList == nil {
		return nil, errors.New("Insertable{{ $t.GoName }}Model is nil")
	}

	valueSetList := make([]sqli.ValuesSetSt, len(modelsList))

	for i, model := range modelsList {
		if model == nil {
			return nil, errors.New("InsertableUserModel is nil")
		}

		valueSetList[i] = sqli.ValueSet(
			{{ range $t.Fields }}
				{{- if not .IsSequence -}}
					sqli.VALUE({{ $t.GoName }}Table.{{ .GoName }}, model.{{ .GoName }}),
				{{ end -}}
			{{ end }}
		)
	}

	query, err := sqli.Query(
		sqli.INSERT_INTO(
			{{ $t.GoName }}Table,
			{{ range $t.Fields }}
				{{- if not .IsSequence -}}
					{{ $t.GoName }}Table.{{ .GoName }},
				{{ end -}}
			{{ end }}
		),
		sqli.VALUES(
			valueSetList...
		),
	)
	if err != nil {
		return nil, err
	}

	return db.ExecContext(ctx, query.SQL, query.Args...)
}

func InsertInto{{ $t.GoName }}TableReturningAll(
	ctx context.Context,
	db DB,
	modelsList ...*Insertable{{ $t.GoName }}Model,
) (*{{ $t.GoName }}Model, error) {
	if modelsList == nil {
		return nil, errors.New("Insertable{{ $t.GoName }}Model is nil")
	}

	valueSetList := make([]sqli.ValuesSetSt, len(modelsList))

	for i, model := range modelsList {
		if model == nil {
			return nil, errors.New("InsertableUserModel is nil")
		}

		valueSetList[i] = sqli.ValueSet(
			{{- range $t.Fields }}
				{{- if not .IsSequence -}}
					{{- if eq .Type "[]uuid.UUID" }}
						pq.Array(model.{{ .GoName }}),
					{{ else }}
						sqli.VALUE({{ $t.GoName }}Table.{{ .GoName }}, model.{{ .GoName }}),
					{{- end -}}
				{{ end -}}
			{{ end }}
		)
	}

	query, err := sqli.Query(
		sqli.INSERT_INTO(
			{{ $t.GoName }}Table,
			{{ range $t.Fields }}
				{{- if not .IsSequence -}}
					{{ $t.GoName }}Table.{{ .GoName }},
				{{ end -}}
			{{ end }}
		),
		sqli.VALUES(
			valueSetList...
		),
		sqli.RETURNING({{ $t.GoName }}Table.AllColumns()),
	)
	if err != nil {
		return nil, err
	}

	row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
	var model {{ $t.GoName }}Model
	err = row.Scan(
		{{- range $t.Fields }}
			{{- if eq .Type "[]uuid.UUID" }}
				pq.Array(&model.{{ .GoName }}),
			{{ else }}
				&model.{{ .GoName }},
			{{- end -}}
		{{ end }}
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// ## Updatable

type Updatable{{ $t.GoName }}Model struct {
	{{ range $t.Fields -}}
		{{ .GoName }} *{{ .Type }} `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
	{{ end }}
}

func NewUpdatable{{ $t.GoName }}Model(
	{{ range $t.Fields -}}
		{{ .GoName }} *{{ .Type }},
	{{ end }}
) *Updatable{{ $t.GoName }}Model {
	return &Updatable{{ $t.GoName }}Model{
		{{ range $t.Fields -}}
			{{ .GoName }},
		{{ end }}
	}
}

{{ range $t.Fields }}
	{{ if (and .IsSequence (not .IsPrimary)) }}
		// ## Select one by sequence {{ .GoName }}
		func Select{{ $t.GoName }}TableBy{{ .GoName }}(
			ctx context.Context,
			db DB,
			{{ .GoName }} {{ .Type }},
		) (*{{ $t.GoName }}Model, error) {
			query, err := sqli.Query(
				sqli.SELECT(
					{{ $t.GoName }}Table.AllColumns(),
				),
				sqli.FROM({{ $t.GoName }}Table),
				sqli.WHERE(
					sqli.EQUAL({{ $t.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
				),
				sqli.LIMIT(1),
			)
			if err != nil {
				return nil, err
			}

			row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
			model := &{{ $t.GoName }}Model{}
			err = row.Scan(
				{{- range $t.Fields }}
					{{- if eq .Type "[]uuid.UUID" }}
						pq.Array(&model.{{ .GoName }}),
					{{ else }}
						&model.{{ .GoName }},
					{{- end -}}
				{{ end }}
			)
			if err != nil {
				return nil, err
			}

			return model, nil
		}

		// ## Delete by sequence {{ .GoName }}
		func DeleteFrom{{ $t.GoName }}TableBy{{ .GoName }}(
			ctx context.Context,
			db DB,
			{{ .GoName }} {{ .Type }},
		) (sql.Result, error) {
			query, err := sqli.Query(
				sqli.DELETE_FROM(
					{{ $t.GoName }}Table,
				),
				sqli.WHERE(
					sqli.EQUAL({{ $t.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
				),
			)
			if err != nil {
				return nil, err
			}

			return db.ExecContext(ctx, query.SQL, query.Args...)
		}
	{{ end }}
{{- end }}

{{ end }}

{{ define "index" }}
{{- $i := .Data -}}

{{ if or $i.IsUnique $i.IsPrimary }}
	{{ if not (eq (len $i.Fields) 1) }}
		// ## Select by compound
		func SelectFrom{{ $i.Table.GoName }}TableBy{{ range $i.Fields }}{{ .GoName }}{{ end }}(
			ctx context.Context,
			db DB,
			{{- range $i.Fields }}
				{{ .GoName }} {{ .Type }},
			{{- end -}}
		) (*{{ $i.Table.GoName }}Model, error) {
			query, err := sqli.Query(
				sqli.SELECT(
					{{ $i.Table.GoName }}Table.AllColumns(),
				),
				sqli.FROM({{ $i.Table.GoName }}Table),
				sqli.WHERE(
					sqli.AND(
						{{ range $i.Fields -}}
							sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
						{{ end }}
					),
				),
				sqli.LIMIT(1),
			)
			if err != nil {
				return nil, err
			}

			row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
			model := &{{ $i.Table.GoName }}Model{}
			err = row.Scan(
				{{- range $i.Table.Fields }}
					{{- if eq .Type "[]uuid.UUID" }}
						pq.Array(&model.{{ .GoName }}),
					{{ else }}
						&model.{{ .GoName }},
					{{- end -}}
				{{ end }}
			)
			if err != nil {
				return nil, err
			}

			return model, nil
		}

		// ## Delete by compound
		func DeleteFrom{{ $i.Table.GoName }}TableBy{{ range $i.Fields }}{{ .GoName }}{{ end }}(
			ctx context.Context,
			db DB,
			{{- range $i.Fields }}
				{{ .GoName }} {{ .Type }},
			{{- end -}}
		) (sql.Result, error) {
			query, err := sqli.Query(
				sqli.DELETE_FROM(
					{{ $i.Table.GoName }}Table,
				),
				sqli.WHERE(
					sqli.AND(
						{{ range $i.Fields -}}
							sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
						{{ end }}
					),
				),
			)
			if err != nil {
				return nil, err
			}

			return db.ExecContext(ctx, query.SQL, query.Args...)
		}

		// ## Update by compound
		func Update{{ $i.Table.GoName }}TableBy{{ range $i.Fields }}{{ .GoName }}{{ end }}(
			ctx context.Context,
			db DB,
			{{- range $i.Fields }}
				{{ .GoName }} {{ .Type }},
			{{- end }}
			updatableModel *Updatable{{ $i.Table.GoName }}Model,
		) (sql.Result, error) {
			valuesSetList := []sqli.Statement{}

			{{ range $i.Table.Fields -}}
				if updatableModel.{{ .GoName }} != nil {
					valuesSetList = append(valuesSetList, sqli.SET_VALUE({{ $i.Table.GoName }}Table.{{ .GoName }}, *updatableModel.{{ .GoName }}))
				}
			{{ end }}

			query, err := sqli.Query(
				sqli.UPDATE(
					{{ $i.Table.GoName }}Table,
				),
				sqli.SET(
					valuesSetList...,
				),
				sqli.WHERE(
					sqli.AND(
						{{ range $i.Fields -}}
							sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
						{{ end }}
					),
				),
			)
			if err != nil {
				return nil, err
			}

			return db.ExecContext(ctx, query.SQL, query.Args...)
		}

		type InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}Result struct {
			{{- range $i.Fields }}
				{{ .GoName }} {{ .Type }} `json:"{{ .SQLName }}" db:"{{ .SQLName }}"`
			{{- end -}}
		}

		func InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}(
			ctx context.Context,
			db DB,
			modelsList ...*Insertable{{ $i.Table.GoName }}Model,
		) (*InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}Result, error) {
			if modelsList == nil {
				return nil, errors.New("InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}Result is nil")
			}

			valueSetList := make([]sqli.ValuesSetSt, len(modelsList))

			for i, model := range modelsList {
				if model == nil {
					return nil, errors.New("InsertableUserModel is nil")
				}

				valueSetList[i] = sqli.ValueSet(
					{{ range $i.Table.Fields }}
						{{- if not .IsSequence -}}
							sqli.VALUE({{ $i.Table.GoName }}Table.{{ .GoName }}, model.{{ .GoName }}),
						{{ end -}}
					{{ end }}
				)
			}

			query, err := sqli.Query(
				sqli.INSERT_INTO(
					{{ $i.Table.GoName }}Table,
					{{ range $i.Table.Fields }}
						{{- if not .IsSequence -}}
							{{ $i.Table.GoName }}Table.{{ .GoName }},
						{{ end -}}
					{{ end }}
				),
				sqli.VALUES(
					valueSetList...
				),
				sqli.RETURNING(
					{{- range $i.Fields }}
						{{ $i.Table.GoName }}Table.{{ .GoName }},
					{{- end }}
				),
			)
			if err != nil {
				return nil, err
			}

			row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
			returning := &InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}Result{}
			err = row.Scan(returning);
			if err != nil {
				return nil, err
			}

			return returning, nil
		}

		{{ range $i.Fields -}}
			// ## Select by {{ .GoName }}
			func Select{{ $i.Table.GoName }}TableBy{{ .GoName }}(
				ctx context.Context,
				db DB,
				{{ .GoName }} {{ .Type }},
			) ([]*{{ $i.Table.GoName }}Model, error) {
				query, err := sqli.Query(
					sqli.SELECT(
						{{ $i.Table.GoName }}Table.AllColumns(),
					),
					sqli.FROM({{ $i.Table.GoName }}Table),
					sqli.WHERE(
						sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
					),
				)
				if err != nil {
					return nil, err
				}

				rows, err := db.QueryxContext(ctx, query.SQL, query.Args...)
				if err != nil {
					return nil, err
				}
				defer rows.Close()
				var items []*{{ $i.Table.GoName }}Model
				for rows.Next() {
					item := &{{ $i.Table.GoName }}Model{}
					if err := rows.Scan(
						item,
					); err != nil {
						return nil, err
					}
					items = append(items, item)
				}
				if err := rows.Err(); err != nil {
					return nil, err
				}
				return items, nil
			}
		{{ end }}
	{{ else }}
		{{ range $i.Fields -}}
			// ## Select by {{ .GoName }}
			func Select{{ $i.Table.GoName }}TableBy{{ .GoName }}(
				ctx context.Context,
				db DB,
				{{ .GoName }} {{ .Type }},
			) (*{{ $i.Table.GoName }}Model, error) {
				query, err := sqli.Query(
					sqli.SELECT(
						{{ $i.Table.GoName }}Table.AllColumns(),
					),
					sqli.FROM({{ $i.Table.GoName }}Table),
					sqli.WHERE(
						sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
					),
					sqli.LIMIT(1),
				)
				if err != nil {
					return nil, err
				}

				row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
				model := &{{ $i.Table.GoName }}Model{}
				err = row.Scan(
					{{- range $i.Table.Fields }}
						{{- if eq .Type "[]uuid.UUID" }}
							pq.Array(&model.{{ .GoName }}),
						{{ else }}
							&model.{{ .GoName }},
						{{- end -}}
					{{ end }}
				)
				if err != nil {
					return nil, err
				}

				return model, nil
			}
		{{ end }}
	{{ end }}
	{{ range $i.Fields -}}
		// ## Delete by {{ .GoName }}
		func DeleteFrom{{ $i.Table.GoName }}TableBy{{ .GoName }}(
			ctx context.Context,
			db DB,
			{{ .GoName }} {{ .Type }},
		) (sql.Result, error) {
			query, err := sqli.Query(
				sqli.DELETE_FROM(
					{{ $i.Table.GoName }}Table,
				),
				sqli.WHERE(
					sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
				),
			)
			if err != nil {
				return nil, err
			}

			return db.ExecContext(ctx, query.SQL, query.Args...)
		}

		func InsertInto{{ $i.Table.GoName }}TableReturning{{ .GoName }}(
			ctx context.Context,
			db DB,
			modelsList ...*Insertable{{ $i.Table.GoName }}Model,
		) (*{{ .Type }}, error) {
			if modelsList == nil {
				return nil, errors.New("InsertInto{{ $i.Table.GoName }}TableReturning{{ range $i.Fields }}{{ .GoName }}{{ end }}Result is nil")
			}

			valueSetList := make([]sqli.ValuesSetSt, len(modelsList))

			for i, model := range modelsList {
				if model == nil {
					return nil, errors.New("InsertableUserModel is nil")
				}

				valueSetList[i] = sqli.ValueSet(
					{{ range $i.Table.Fields }}
						{{- if not .IsSequence -}}
							sqli.VALUE({{ $i.Table.GoName }}Table.{{ .GoName }}, model.{{ .GoName }}),
						{{ end -}}
					{{ end }}
				)
			}

			query, err := sqli.Query(
				sqli.INSERT_INTO(
					{{ $i.Table.GoName }}Table,
					{{ range $i.Table.Fields }}
						{{- if not .IsSequence -}}
							{{ $i.Table.GoName }}Table.{{ .GoName }},
						{{ end -}}
					{{ end }}
				),
				sqli.VALUES(
					valueSetList...
				),
				sqli.RETURNING(
					{{ $i.Table.GoName }}Table.{{ .GoName }},
				),
			)
			if err != nil {
				return nil, err
			}

			row := db.QueryRowxContext(ctx, query.SQL, query.Args...)
			var returning {{ .Type }}
			err = row.Scan(&returning);
			if err != nil {
				return nil, err
			}

			return &returning, nil
		}

		// # Update
		// ## Update by {{ .GoName }}
		func Update{{ $i.Table.GoName }}TableBy{{ .GoName }}(
			ctx context.Context,
			db DB,
			{{ .GoName }} {{ .Type }},
			updatableModel *Updatable{{ $i.Table.GoName }}Model,
		) (sql.Result, error) {
			valuesSetList := []sqli.Statement{}

			{{ range $i.Table.Fields -}}
				if updatableModel.{{ .GoName }} != nil {
					valuesSetList = append(valuesSetList, sqli.SET_VALUE({{ $i.Table.GoName }}Table.{{ .GoName }}, *updatableModel.{{ .GoName }}))
				}
			{{ end }}

			query, err := sqli.Query(
				sqli.UPDATE(
					{{ $i.Table.GoName }}Table,
				),
				sqli.SET(
					valuesSetList...,
				),
				sqli.WHERE(
					sqli.EQUAL({{ $i.Table.GoName }}Table.{{ .GoName }}, {{ .GoName }}),
				),
			)
			if err != nil {
				return nil, err
			}

			return db.ExecContext(ctx, query.SQL, query.Args...)
		}
	{{ end }}
{{ end }}

{{end}}