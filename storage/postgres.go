package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

func InitPostgresDB(env map[string]string) *PostgresDB {
	res := &PostgresDB{
		env: env,
	}

	res.NewPsqlDB()

	return res
}

type PostgresDB struct {
	DB     *sql.DB
	env    map[string]string
	cancel func()
}

// NewPsqlDB создает новое соединение с базой данных PostgreSQL
func (p *PostgresDB) NewPsqlDB() {
	var err error
	err = p.tryConnect()
	if err != nil {
		logrus.Errorf("Postgresql init error: %s", err)
		time.Sleep(5 * time.Second) // Wait before retrying
		err = p.tryConnect()
	}

	err = p.DB.Ping()
	if err != nil {
		logrus.Errorf("Postgresql ping error: %s", err)
		time.Sleep(5 * time.Second) // Wait before retrying
		err = p.DB.Ping()
	}

	logrus.Info("Successful connection to the database")
}

func (p *PostgresDB) tryConnect() error {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		p.env["POSTGRES_HOST"],
		p.env["POSTGRES_PORT"],
		p.env["POSTGRES_DEFAULT_USER"],
		p.env["POSTGRES_DBNAME"],
		p.env["POSTGRES_DEFAULT_PASS"],
	)
	var err error
	p.DB, err = sql.Open(p.env["PG_DRIVER"], dataSourceName)
	if err != nil {
		return err
	}
	return nil
}

// InsertMessage вставляет новое сообщение в таблицу сообщений и возвращает идентификатор нового сообщения.
func (p *PostgresDB) InsertMessage(messageBody string) (int, error) {
	var messageID int
	query := `INSERT INTO message_schema.messages (message_body) VALUES ($1) RETURNING id`
	err := p.DB.QueryRow(query, messageBody).Scan(&messageID)
	if err != nil {
		return 0, err
	}
	return messageID, nil
}
