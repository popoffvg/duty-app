package service

import (
	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

type TeammateService struct {
	repo repository.Teammate
}

func NewTeammateService(repo repository.Teammate) *TeammateService {
	return &TeammateService{repo: repo}
}

func (ts *TeammateService) Create(userId, teamId int, mate model.Teammate) (int, error) {
	return ts.repo.Create(userId, teamId, mate)
}

func (ts *TeammateService) Read(userId, teammateId int) (model.Teammate, error) {
	return ts.repo.Read(userId, teammateId)
}

func (ts *TeammateService) Update(userId, teammateId int, input model.UpdateTeammateInput) error {
	return ts.repo.Update(userId, teammateId, input)
}

func (ts *TeammateService) Delete(userId, teammateId int) error {
	return ts.repo.Delete(userId, teammateId)
}

func (ts *TeammateService) List(userId, teamId int) ([]model.Teammate, error) {
	return ts.repo.List(userId, teamId)
}
