package service

import (
	"time"

	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

type DutyService struct {
	repo repository.Duty
}

func NewDutyService(repo repository.Duty) *DutyService {
	return &DutyService{repo: repo}
}

func (ds *DutyService) Create(userId, teamId int, duty model.Duty) (model.Duty, error) {
	timeNow := time.Now()

	if duty.IsDaily {
		// daily duty: set current date
		duty.Date = timeNow
	} else {
		// weekly duty: set date of Monday of this week
		duty.Date = ds.getMondayDate(timeNow)
	}

	return ds.repo.Create(userId, teamId, duty)
}

func (ds *DutyService) Update(userId, dutyId int, input model.UpdateDutyInput) error {
	return ds.repo.Update(userId, dutyId, input)
}

func (ds *DutyService) Delete(userId, dutyId int) error {
	return ds.repo.Delete(userId, dutyId)
}

func (ds *DutyService) List(userId, teamId int) ([]model.Duty, error) {
	return ds.repo.List(userId, teamId)
}

// Non-CRUDL methods

func (ds *DutyService) ReadCurrent(userId, teamId int) ([]model.Duty, error) {
	dailyDate := time.Now()
	weeklyDate := ds.getMondayDate(dailyDate)

	return ds.repo.ReadCurrent(userId, teamId, dailyDate, weeklyDate)
}

func (ds *DutyService) History(userId, teamId int) ([]model.History, error) {
	return ds.repo.History(userId, teamId)
}

func (ds *DutyService) getMondayDate(date time.Time) time.Time {
	var mondayDate time.Time

	// weekly duty: try to find date of Monday of this week
	weekdayNumber := int(date.Weekday()) // Sunday=0, Monday=1
	switch weekdayNumber {
	case 0: // Sunday: -6 days to Monday
		mondayDate = date.AddDate(0, 0, -6)
	case 1: // Monday: set this date
		mondayDate = date
	default: // Another weekday: set date of Monday by day diff between
		mondayDate = date.AddDate(0, 0, 1-weekdayNumber)
	}

	return mondayDate
}
