package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vitalygudza/duty-app/internal/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(signInInput model.SignInInput) (model.User, error)
}

type Team interface {
	Create(userId int, team model.Team) (int, error)
	Read(userId, teamId int) (model.Team, error)
	Update(userId, teamId int, input model.UpdateTeamInput) error
	Delete(userId, teamId int) error
	List(userId int) ([]model.Team, error)
}

type Teammate interface {
	Create(userId, teamId int, team model.Teammate) (int, error)
	Read(userId, teammateId int) (model.Teammate, error)
	Update(userId, teammateId int, input model.UpdateTeammateInput) error
	Delete(userId, teammateId int) error
	List(userId, teamId int) ([]model.Teammate, error)
}

type Duty interface {
	Create(userId, teamId int, duty model.Duty) (model.Duty, error)
	Update(userId, dutyId int, input model.UpdateDutyInput) error
	Delete(userId, dutyId int) error
	List(userId, teamId int) ([]model.Duty, error)

	ReadCurrent(userId, teamId int, dailyDate, weeklyDate time.Time) ([]model.Duty, error)
	History(userId, teamId int) ([]model.History, error)
}

type Repository struct {
	Authorization
	Team
	Teammate
	Duty
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Team:          NewTeamPostgres(db),
		Teammate:      NewTeammatePostgres(db),
		Duty:          NewDutyPostgres(db),
	}
}
