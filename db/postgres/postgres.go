package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "12345"
	dbName   = "wb"
	port     = "5433"
	sslMode  = "disable"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbCfg())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func dbCfg() string {
	ans :=
		"host=" + host + " user=" + user +
			" password=" + password +
			" dbname=" + dbName + " port=" + port +
			" sslmode=" + sslMode
	return ans
}
