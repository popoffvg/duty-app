package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vitalygudza/duty-app/internal/model"
)

type DutyPostgres struct {
	db *sqlx.DB
}

func NewDutyPostgres(db *sqlx.DB) *DutyPostgres {
	return &DutyPostgres{db: db}
}

// todo refactoring - why transaction?
func (dp *DutyPostgres) Create(userId, teamId int, duty model.Duty) (model.Duty, error) {
	tx, err := dp.db.Begin()
	if err != nil {
		return duty, err
	}

	// check: can user create duty in team? Is user is owner of team?
	var id int
	checkUserTeamQuery := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1 AND id = $2", teamsTable)
	row := tx.QueryRow(checkUserTeamQuery, userId, teamId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return duty, err
	}

	createDutyQuery := fmt.Sprintf(
		"INSERT INTO %s (team_id, teammate_id, is_daily, date) VALUES ($1, $2, $3, $4) RETURNING id",
		dutiesTable)

	row = tx.QueryRow(createDutyQuery, teamId, duty.TeammateId, duty.IsDaily, duty.Date)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return duty, err
	}

	duty.Id = id
	return duty, tx.Commit()
}

func (dp *DutyPostgres) Update(userId, dutyId int, input model.UpdateDutyInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.TeammateId != nil {
		setValues = append(setValues, fmt.Sprintf("teammate_id=$%d", argId))
		args = append(args, *input.TeammateId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s dt SET %s FROM %s t 
				WHERE dt.team_id=t.id AND t.user_id = $%d AND dt.id = $%d`,
		dutiesTable, setQuery, teamsTable, argId, argId+1)
	args = append(args, userId, dutyId)

	_, err := dp.db.Exec(query, args...)
	return err
}

func (dp *DutyPostgres) Delete(userId, dutyId int) error {
	query := fmt.Sprintf(`DELETE FROM %s dt USING %s t WHERE dt.team_id=t.id AND 
									t.user_id = $1 AND dt.id = $2`, dutiesTable, teamsTable)
	_, err := dp.db.Exec(query, userId, dutyId)
	return err
}

func (dp *DutyPostgres) List(userId, teamId int) ([]model.Duty, error) {
	var duties []model.Duty
	query := fmt.Sprintf(`SELECT dt.id, dt.team_id, dt.teammate_id, dt.is_daily, dt.date 
									FROM %s dt INNER JOIN %s t on dt.team_id=t.id
									WHERE t.user_id = $1 AND dt.team_id = $2
									ORDER BY dt.teammate_id`, dutiesTable, teamsTable)
	if err := dp.db.Select(&duties, query, userId, teamId); err != nil {
		return nil, err
	}

	return duties, nil
}

func (dp *DutyPostgres) ReadCurrent(userId, teamId int, dailyDate, weeklyDate time.Time) ([]model.Duty, error) {
	var duties []model.Duty

	query := fmt.Sprintf(`SELECT dt.id, dt.team_id, dt.teammate_id, dt.is_daily, dt.date 
									FROM %s dt INNER JOIN %s t on dt.team_id=t.id
									WHERE t.user_id = $1 AND dt.team_id = $2
									AND (dt.is_daily = true AND dt.date = $3 OR
										dt.is_daily = false AND dt.date = $4)`,
		dutiesTable, teamsTable)
	if err := dp.db.Select(&duties, query, userId, teamId, dailyDate, weeklyDate); err != nil {
		return nil, err
	}

	return duties, nil
}

func (dp *DutyPostgres) History(userId, teamId int) ([]model.History, error) {
	var history []model.History
	query := fmt.Sprintf(`SELECT dt.id, tms.name, dt.date, dt.is_daily FROM %s dt 
									INNER JOIN %s t on dt.team_id=t.id
									INNER JOIN %s tms on dt.teammate_id=tms.id
									WHERE t.user_id = $1 AND dt.team_id = $2
									ORDER BY dt.date DESC LIMIT 100`, dutiesTable, teamsTable, teammatesTable)
	if err := dp.db.Select(&history, query, userId, teamId); err != nil {
		return nil, err
	}

	return history, nil
}

// todo refactoring: привести к единообразию в sql-запросах - либо с алиасами, либо без них (...FROM table AS asd)

func (dp *DutyPostgres) FindLastDuty(isDaily bool, date time.Time) ([]model.Duty, error) {
	var duties []model.Duty

	// get top 1 (in group)
	// select distinct on (team_id) id, team_id, date from duties where date > current_date - 6 order by team_id, date desc;

	// get top N (in group) - window function https://ubiq.co/database-blog/how-to-get-first-row-per-group-in-postgresql/
	// select * from (select *, row_number() over (partition by team_id order by date desc) as row_number from duties where date > current_date - 6) temp where row_number=N;

	// get last daily duty for every team in 3 last days
	query := fmt.Sprintf(`SELECT DISTINCT ON (team_id)
									id, team_id, teammate_id, is_daily, date 
									FROM %s WHERE is_daily = $1 AND date > $2
									ORDER BY team_id, date desc`, dutiesTable)
	if err := dp.db.Select(&duties, query, isDaily, date); err != nil {
		return nil, err
	}

	return duties, nil
}

func (dp *DutyPostgres) CreateNextDuty(duty model.Duty) (int, error) {
	var id int
	var exists bool

	// check duty intersection (with duty.date duty)
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s
									WHERE team_id=$1 AND is_daily=$2 AND date=$3)`, dutiesTable)
	row := dp.db.QueryRow(query, duty.TeamId, duty.IsDaily, duty.Date)
	if err := row.Scan(&exists); err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("duty on date %q is already exists (is daily?: %t)", duty.Date, duty.IsDaily)
	}

	// create new duty
	query = fmt.Sprintf(`INSERT INTO %s (team_id, teammate_id, is_daily, date) 
									VALUES ($1, $2, $3, $4) RETURNING id`, dutiesTable)
	row = dp.db.QueryRow(query, duty.TeamId, duty.TeammateId, duty.IsDaily, duty.Date)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
