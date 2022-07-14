package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vitalygudza/duty-app/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (ap *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int
	var exists bool

	// check if username already exists
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE username=$1)", usersTable)
	row := ap.db.QueryRow(query, user.Username)
	if err := row.Scan(&exists); err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("username %q is already exists", user.Username)
	}

	query = fmt.Sprintf("INSERT INTO %s (name, username, password) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row = ap.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (ap *AuthPostgres) GetUser(signInInput model.SignInInput) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password=$2", usersTable)
	err := ap.db.Get(&user, query, signInInput.Username, signInInput.Password)

	return user, err
}
