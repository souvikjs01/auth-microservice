package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/souvikjs01/auth-microservice/config"
)

const (
	maxOpenConns    = 60
	connMaxLifeTime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dsn := c.Postgres.PgURI

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifeTime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
