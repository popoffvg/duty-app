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
	query := fmt.Sprintf("INSERT INTO %s (title, user_id, space_channel) VALUES ($1, $2, $3) RETURNING id", teamsTable)
	row := tp.db.QueryRow(query, team.Title, userId, team.SpaceChannel)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tp *TeamPostgres) List(userId int) ([]model.Team, error) {
	var teams []model.Team
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY id", teamsTable)
	err := tp.db.Select(&teams, query, userId)

	return teams, err
}

func (tp *TeamPostgres) Read(userId, teamId int) (model.Team, error) {
	var team model.Team
	// todo: select * --> select row1, row2, ...
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	err := tp.db.Get(&team, query, userId, teamId)

	return team, err
}

func (tp *TeamPostgres) Update(userId int, teamId int, input model.UpdateTeamInput) error {
	// todo: update if row is exists!!
	query := fmt.Sprintf("UPDATE %s SET title = $1, space_channel = $2 WHERE user_id = $3 AND id = $4", teamsTable)
	_, err := tp.db.Exec(query, input.Title, input.SpaceChannel, userId, teamId)

	return err
}

func (tp *TeamPostgres) Delete(userId, teamId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	_, err := tp.db.Exec(query, userId, teamId)

	return err
}

func (tp *TeamPostgres) GetTeamInfo(teamId int) (model.Team, error) {
	var team model.Team
	query := fmt.Sprintf("SELECT title FROM %s WHERE id = $1", teamsTable)
	err := tp.db.Get(&team, query, teamId)

	return team, err
}
