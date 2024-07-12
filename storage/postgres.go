package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

// InitPostgresDB инициализирует новый экземпляр PostgresDB с заданными переменными среды.
// Он устанавливает новое соединение с базой данных, используя предоставленный контекст для отмены.
// Возвращает указатель на инициализированный экземпляр PostgresDB.
func InitPostgresDB(env map[string]string) *PostgresDB {
	ctx, cancel := context.WithCancel(context.Background())
	res := &PostgresDB{
		env:    env,
		cancel: cancel,
	}

	go res.ConnectPostgres(ctx)
	return res
}

// PostgresDB содержит соединение с базой данных, переменные среды и функцию отмены.
// для прекращения попытки подключения к базе данных.
type PostgresDB struct {
	DB     *sql.DB
	env    map[string]string
	cancel func()
}

// ConnectPostgres создает новое соединение с базой данных PostgreSQL
func (p *PostgresDB) ConnectPostgres(ctx context.Context) {
	var err error
	err = p.tryConnect(ctx)
	for err != nil {
		select {
		case <-ctx.Done():
			logrus.Errorf("Context cancelled, stopping connection attempts")
			return
		default:
			logrus.Errorf("Postgresql init error: %s", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			err = p.tryConnect(ctx)
		}
	}
	logrus.Info("Successful connection to the database")
}

// tryConnect пытается подключиться к базе данных PostgreSQL, используя конфигурацию
// указанный в структуре PostgresDB. Он использует предоставленный контекст, чтобы разрешить отмену.
// попытки подключения. Возвращает ошибку, если попытка подключения не удалась.
func (p *PostgresDB) tryConnect(ctx context.Context) error {
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
	return p.DB.PingContext(ctx)
}

// RecordMessage записывает новое сообщение в таблицу сообщений и возвращает идентификатор нового сообщения.
func (p *PostgresDB) RecordMessage(messageBody []byte) (int, error) {
	var messageID int
	query := `INSERT INTO message_schema.messages (message_body) VALUES ($1) RETURNING id`
	err := p.DB.QueryRow(query, messageBody).Scan(&messageID)
	if err != nil {
		return 0, err
	}
	return messageID, nil
}
