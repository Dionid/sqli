package db

type TablesSt struct {
	GooseDbVersion string `json:"goose_db_version" db:"goose_db_version"`
	Office         string `json:"office" db:"office"`
	OfficeUser     string `json:"office_user" db:"office_user"`
	User           string `json:"user" db:"user"`
}

var Tables = TablesSt{
	GooseDbVersion: "goose_db_version",
	Office:         "office",
	OfficeUser:     "office_user",
	User:           "user",
}

// Named "T" for shortness
var T = Tables
