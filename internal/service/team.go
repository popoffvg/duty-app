package service

import (
	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

type TeamService struct {
	repo repository.Team
}

func NewTeamService(repo repository.Team) *TeamService {
	return &TeamService{repo: repo}
}

func (ts *TeamService) Create(userId int, team model.Team) (int, error) {
	return ts.repo.Create(userId, team)
}

func (ts *TeamService) Read(userId int, teamId int) (model.Team, error) {
	return ts.repo.Read(userId, teamId)
}

func (ts *TeamService) Update(userId int, teamId int, input model.UpdateTeamInput) error {
	return ts.repo.Update(userId, teamId, input)
}

func (ts *TeamService) Delete(userId int, teamId int) error {
	return ts.repo.Delete(userId, teamId)
}

func (ts *TeamService) List(userId int) ([]model.Team, error) {
	return ts.repo.List(userId)
}
