package store

import (
	"errors"
	"marketplace/internal/model"
	"marketplace/pkg/hash"
)

func (db *DB) CreateUser(u *model.User) error {
	if db.QueryRow("SELECT id FROM users WHERE username = $1", u.Username).Scan(&u.ID) == nil {
		return errors.New("User already exists")
	}
	hashedPassword, err := hash.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return db.QueryRow("INSERT INTO users(username, password) VALUES($1, $2) RETURNING id", u.Username, u.Password).Scan(&u.ID)
}

func (db *DB) GetUser(username string) (*model.User, error) {
	u := &model.User{}
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}
