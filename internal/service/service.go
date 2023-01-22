package service

import (
	"time"

	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	CheckCredentials(signInInput model.SignInInput) (model.User, error)
	GenerateToken(expireTime time.Time, userId int) (string, error)
	ParseToken(token string) (int, error)
}

type Team interface {
	Create(userId int, team model.Team) (int, error)
	Read(userId int, teamId int) (model.Team, error)
	Update(userId int, teamId int, input model.UpdateTeamInput) error
	Delete(userId int, teamId int) error
	List(userId int) ([]model.Team, error)

	GetTeamInfo(teamId int) (model.Team, error)
}

type Teammate interface {
	Create(userId, teamId int, input model.Teammate) (int, error)
	Update(userId, teammateId int, input model.UpdateTeammateInput) error
	Read(userId, teammateId int) (model.Teammate, error)
	Delete(userId, teammateId int) error
	List(userId, teamId int) ([]model.Teammate, error)

	ListReadyTeammates(teamId int) ([]model.Teammate, error)
}

type Duty interface {
	Create(userId, teamId int, input model.Duty) (model.Duty, error)
	Update(userId, dutyId int, input model.UpdateDutyInput) error
	Delete(userId, dutyId int) error
	List(userId, teamId int) ([]model.Duty, error)

	ReadCurrent(userId, teamId int) ([]model.Duty, error)
	History(userId, teamId int) ([]model.History, error)

	UpdateDutiesJob(timeNow time.Time) ([]model.Notification, error)
}

type Notifier interface {
	SendNotification(channelName, text string) error
}

type Service struct {
	Authorization
	Team
	Teammate
	Duty
	Notifier
}

func NewService(repo *repository.Repository, notifier Notifier) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Team:          NewTeamService(repo.Team),
		Teammate:      NewTeammateService(repo.Teammate),
		Duty:          NewDutyService(repo.Team, repo.Teammate, repo.Duty),
		Notifier:      notifier,
	}
}
