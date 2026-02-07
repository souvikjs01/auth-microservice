package postgres

import (
	"fmt"
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
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Postgres.Host, c.Server.Port, c.Postgres.User, c.Postgres.Dbname, c.Postgres.Password)

	db, err := sqlx.Connect(c.Postgres.PgDriver, dsn)
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
