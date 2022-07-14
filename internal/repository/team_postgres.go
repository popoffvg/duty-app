package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vitalygudza/duty-app/internal/model"
)

type TeamPostgres struct {
	db *sqlx.DB
}

func NewTeamPostgres(db *sqlx.DB) *TeamPostgres {
	return &TeamPostgres{db: db}
}

func (tp *TeamPostgres) Create(userId int, team model.Team) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (title, user_id) VALUES ($1, $2) RETURNING id", teamsTable)
	row := tp.db.QueryRow(query, team.Title, userId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tp *TeamPostgres) List(userId int) ([]model.Team, error) {
	var teams []model.Team
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", teamsTable)
	err := tp.db.Select(&teams, query, userId)

	return teams, err
}

func (tp *TeamPostgres) Read(userId, teamId int) (model.Team, error) {
	var team model.Team
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	err := tp.db.Get(&team, query, userId, teamId)

	return team, err
}
func (tp *TeamPostgres) Update(userId int, teamId int, input model.UpdateTeamInput) error {
	// todo: update if row is exists!!
	query := fmt.Sprintf("UPDATE %s SET title = $1 WHERE user_id = $2 AND id = $3", teamsTable)
	_, err := tp.db.Exec(query, input.Title, userId, teamId)

	return err
}

func (tp *TeamPostgres) Delete(userId, teamId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	_, err := tp.db.Exec(query, userId, teamId)

	return err
}
