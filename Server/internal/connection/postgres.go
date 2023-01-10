package connection

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func PSQLConnect() (*sql.DB, error) {
	connStr := "user=postgres password=8403 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	return db, err
}
