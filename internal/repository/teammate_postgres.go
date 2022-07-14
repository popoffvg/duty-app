package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vitalygudza/duty-app/internal/model"
)

type TeammatePostgres struct {
	db *sqlx.DB
}

func NewTeammatePostgres(db *sqlx.DB) *TeammatePostgres {
	return &TeammatePostgres{db: db}
}

func (tp *TeammatePostgres) Create(userId, teamId int, mate model.Teammate) (int, error) {
	tx, err := tp.db.Begin()
	if err != nil {
		return 0, err
	}

	// check: can user create teammate in team? Is user is owner of team?
	var id int
	checkUserTeamQuery := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	row := tx.QueryRow(checkUserTeamQuery, userId, teamId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createTeammateQuery := fmt.Sprintf(
		"INSERT INTO %s (team_id, name, duty_readiness, duties) VALUES ($1, $2, $3, $4) RETURNING id",
		teammatesTable)
	row = tx.QueryRow(createTeammateQuery, teamId, mate.Name, mate.DutyReadiness, mate.Duties)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (tp *TeammatePostgres) Read(userId, teammateId int) (model.Teammate, error) {
	var teammate model.Teammate
	query := fmt.Sprintf(`SELECT tms.id, tms.team_id, tms.name, tms.duty_readiness, tms.duties 
									FROM %s tms INNER JOIN %s t on tms.team_id=t.id
									WHERE t.user_id = $1 AND tms.id = $2`, teammatesTable, teamsTable)
	if err := tp.db.Get(&teammate, query, userId, teammateId); err != nil {
		return teammate, err
	}

	return teammate, nil
}

func (tp *TeammatePostgres) Update(userId, teammateId int, input model.UpdateTeammateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.DutyReadiness != nil {
		setValues = append(setValues, fmt.Sprintf("duty_readiness=$%d", argId))
		args = append(args, *input.DutyReadiness)
		argId++
	}

	if input.Duties != nil {
		setValues = append(setValues, fmt.Sprintf("duties=$%d", argId))
		args = append(args, *input.Duties)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s tms SET %s FROM %s t WHERE tms.team_id=t.id AND t.user_id = $%d AND tms.id = $%d`,
		teammatesTable, setQuery, teamsTable, argId, argId+1)
	args = append(args, userId, teammateId)

	_, err := tp.db.Exec(query, args...)
	return err
}

func (tp *TeammatePostgres) Delete(userId, teammateId int) error {
	query := fmt.Sprintf(`DELETE FROM %s tms USING %s t WHERE tms.team_id=t.id AND 
									t.user_id = $1 AND tms.id = $2`, teammatesTable, teamsTable)
	_, err := tp.db.Exec(query, userId, teammateId)
	return err
}

func (tp *TeammatePostgres) List(userId, teamId int) ([]model.Teammate, error) {
	var teammates []model.Teammate
	query := fmt.Sprintf(`SELECT tms.id, tms.team_id, tms.name, tms.duty_readiness, tms.duties 
									FROM %s tms INNER JOIN %s t on tms.team_id=t.id
									WHERE t.user_id = $1 AND tms.team_id = $2 ORDER BY tms.id`, teammatesTable, teamsTable)
	if err := tp.db.Select(&teammates, query, userId, teamId); err != nil {
		return nil, err
	}

	return teammates, nil
}
