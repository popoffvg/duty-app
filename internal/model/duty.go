package model

import (
	"errors"
	"time"
)

type Duty struct {
	Id         int       `json:"id,omitempty" db:"id"`
	TeamId     int       `json:"team_id,omitempty" db:"team_id"`
	TeammateId int       `json:"teammate_id" db:"teammate_id" binding:"required"`
	IsDaily    bool      `json:"is_daily" db:"is_daily"`
	Date       time.Time `json:"date,omitempty" db:"date"`
}

type UpdateDutyInput struct {
	TeammateId *int `json:"teammate_id"`
}

func (u UpdateDutyInput) Validate() error {
	if u.TeammateId == nil {
		return errors.New("update duties: body has no values")
	}
	return nil
}

type History struct {
	Id      int       `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	IsDaily bool      `json:"is_daily" db:"is_daily"`
	Date    time.Time `json:"date,omitempty" db:"date"`
}
