package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPsqlDB(env map[string]string) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		env["POSTGRES_HOST"],
		env["POSTGRES_PORT"],
		env["POSTGRES_DEFAULT_USER"],
		env["POSTGRES_DBNAME"],
		env["POSTGRES_DEFAULT_PASS"],
	)

	db, err := sql.Open(env["PG_DRIVER"], dataSourceName)
	if err != nil {
		return nil, err
	}
	logrus.Info("successful connection to the database")

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
