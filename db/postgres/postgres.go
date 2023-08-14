package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	port     = "5432"
	user     = "postgres"
	password = "12345"
	sslMode  = "disable"
	dbName   = "wb"
	host     = "localhost"
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
