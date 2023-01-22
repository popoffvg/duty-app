package service

import (
	"fmt"
	"time"

	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/repository"
)

const (
	// number of previous days to search for a last daily duty
	previousDaysLookUpForDailyDuty = 3
	// number of previous days to search for a last weekly duty
	previousDaysLookUpForWeeklyDuty = 8
)

type DutyService struct {
	repoTeam     repository.Team
	repoTeammate repository.Teammate
	repoDuty     repository.Duty
}

func NewDutyService(repoTeam repository.Team, repoTeammate repository.Teammate, repoDuty repository.Duty) *DutyService {
	return &DutyService{
		repoTeam:     repoTeam,
		repoTeammate: repoTeammate,
		repoDuty:     repoDuty,
	}
}

func (ds *DutyService) Create(userId, teamId int, duty model.Duty) (model.Duty, error) {
	now := time.Now()

	if duty.IsDaily {
		// daily duty: set current date
		duty.Date = now
	} else {
		// weekly duty: set date of Monday of this week
		duty.Date = ds.getMondayDate(now)
	}

	return ds.repoDuty.Create(userId, teamId, duty)
}

func (ds *DutyService) Update(userId, dutyId int, input model.UpdateDutyInput) error {
	return ds.repoDuty.Update(userId, dutyId, input)
}

func (ds *DutyService) Delete(userId, dutyId int) error {
	return ds.repoDuty.Delete(userId, dutyId)
}

func (ds *DutyService) List(userId, teamId int) ([]model.Duty, error) {
	return ds.repoDuty.List(userId, teamId)
}

func (ds *DutyService) ReadCurrent(userId, teamId int) ([]model.Duty, error) {
	dailyDate := time.Now()
	weeklyDate := ds.getMondayDate(dailyDate)

	return ds.repoDuty.ReadCurrent(userId, teamId, dailyDate, weeklyDate)
}

func (ds *DutyService) History(userId, teamId int) ([]model.History, error) {
	return ds.repoDuty.History(userId, teamId)
}

func (ds *DutyService) UpdateDutiesJob(timeNow time.Time) ([]model.Notification, error) {
	var newDuties []model.Notification

	// daily duty: update only for work days
	if !ds.isWeekend(timeNow) {
		nextDailyDuties, err := ds.SetNextDuty(true, timeNow)
		if err != nil {
			return nil, fmt.Errorf("update daily duty error: %w", err)
		}

		newDuties = append(newDuties, nextDailyDuties...)
	}

	// weekly duty: update only on Monday
	if ds.isMonday(timeNow) {
		nextWeeklyDuties, err := ds.SetNextDuty(false, timeNow)
		if err != nil {
			return nil, fmt.Errorf("update daily duty error: %w", err)
		}

		newDuties = append(newDuties, nextWeeklyDuties...)
	}

	return newDuties, nil
}

func (ds *DutyService) SetNextDuty(isDaily bool, now time.Time) ([]model.Notification, error) {
	var nextDuties []model.Notification

	daysLookup := previousDaysLookUpForDailyDuty
	if !isDaily {
		daysLookup = previousDaysLookUpForWeeklyDuty
	}

	// list of last duties for all teams. Lookup for several days back
	lastDuties, err := ds.repoDuty.FindLastDuty(isDaily, now.AddDate(0, 0, -daysLookup))
	if err != nil {
		return nil, fmt.Errorf("find last duty error: %w", err)
	}

	for _, duty := range lastDuties {
		// get ordered (by id) list of teammates (they are ready to duties)
		teammates, err := ds.repoTeammate.ListReadyTeammates(duty.TeamId)
		if err != nil {
			return nil, fmt.Errorf("can not get teammates list: %w", err)
		}

		if len(teammates) > 0 {
			var nextDutyTeammate *model.Teammate
			for _, teammate := range teammates {
				// get the nearest (next) ready teammate
				if teammate.Id > duty.TeammateId {
					nextDutyTeammate = &teammate
					break
				}
			}

			// if duties is ended, let's start duties soon for first teammate in sorted teammates list
			if nextDutyTeammate == nil {
				nextDutyTeammate = &teammates[0]
			}

			// set new duty
			_, err := ds.repoDuty.CreateNextDuty(model.Duty{
				TeamId:     duty.TeamId,
				TeammateId: nextDutyTeammate.Id,
				IsDaily:    isDaily,
				Date:       now,
			})
			if err != nil {
				return nil, fmt.Errorf("can not create next duty: %w", err)
			}

			// todo get channel of team and name. Now team.Title == team.channelSpace
			team, err := ds.repoTeam.GetTeamInfo(duty.TeamId)
			if err != nil {
				return nil, fmt.Errorf("can not get team info: %w", err)
			}

			// save info about new duty for notifications
			nextDuties = append(nextDuties, model.Notification{
				TeamName:        team.Title,
				TeamChannelName: team.Title,
				TeammateName:    nextDutyTeammate.Name,
				IsDaily:         isDaily,
				Date:            now,
			})
		}
	}

	return nextDuties, nil
}

func (ds *DutyService) isWeekend(date time.Time) bool {
	// Sunday=0, Monday=1
	weekdayNumber := int(date.Weekday())
	return weekdayNumber == 0 || weekdayNumber == 6
}

func (ds *DutyService) isMonday(date time.Time) bool {
	weekdayNumber := int(date.Weekday())
	return weekdayNumber == 1
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
