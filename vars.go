package main

import (
	"database/sql"
	"sync"
)

type job struct {
	query string
	vals  []interface{}
}

var (
	sslMode, dbType, host, port, user, password, dbname string
	db                                                  *sql.DB
	err                                                 error
	wg                                                  sync.WaitGroup
)

var validDBChoices = map[int]string{
	1: "mysql",
	2: "postgres",
	3: "sqlite",
}

var validMysqlChoices = map[int]string{
	1: "VARCHAR(60)",
	2: "INT",
	3: "FLOAT",
	4: "DECIMAL",
	5: "DATE",
	6: "TIME",
}
var validPostgresChoices = map[int]string{
	1: "varchar(60)",
	2: "integer",
	3: "float(12)",
	4: "date",
	5: "time(6)",
}
var validSqliteChoices = map[int]string{
	1: "TEXT",
	2: "INTEGER",
	3: "DECIMAL",
}
var dbTypeChoices = map[string]map[int]string{
	"mysql":    validMysqlChoices,
	"postgres": validPostgresChoices,
	"sqlite":   validSqliteChoices,
}
