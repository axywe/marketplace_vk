package store

import (
	"database/sql"
	"fmt"
	"marketplace/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
	Config *config.Config
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName))
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, Config: cfg}, nil
}
