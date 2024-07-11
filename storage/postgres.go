package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewPsqlDB(env map[string]string) *sql.DB {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		env["POSTGRES_HOST"],
		env["POSTGRES_PORT"],
		env["POSTGRES_DEFAULT_USER"],
		env["POSTGRES_DBNAME"],
		env["POSTGRES_DEFAULT_PASS"],
	)

	db, err := sql.Open(env["PG_DRIVER"], dataSourceName)
	if err != nil {
		logrus.Fatalf("Postgresql init: %s", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatalf("Postgresql init: %s", err)
		return nil
	}
	logrus.Info("successful connection to the database")

	return db
}

// InsertMessage вставляет новое сообщение в таблицу сообщений и возвращает идентификатор нового сообщения.
func InsertMessage(db *sql.DB, messageBody string) (int, error) {
	var messageID int
	query := `INSERT INTO message_schema.messages (message_body) VALUES ($1) RETURNING id`
	err := db.QueryRow(query, messageBody).Scan(&messageID)
	if err != nil {
		return 0, err
	}
	return messageID, nil
}
