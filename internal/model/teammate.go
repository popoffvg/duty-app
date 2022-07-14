package model

import "errors"

type Teammate struct {
	Id            int    `json:"id,omitempty" db:"id"`
	TeamId        int    `json:"team_id,omitempty" db:"team_id"`
	Name          string `json:"name" db:"name" binding:"required,min=3"`
	DutyReadiness bool   `json:"duty_readiness" db:"duty_readiness"`
	Duties        int    `json:"duties" db:"duties"`
}

type UpdateTeammateInput struct {
	Name          *string `json:"name"`
	DutyReadiness *bool   `json:"duty_readiness"`
	Duties        *int    `json:"duties"`
}

func (u UpdateTeammateInput) Validate() error {
	if u.Name == nil && u.DutyReadiness == nil && u.Duties == nil {
		return errors.New("update teammates: body has no values")
	}
	return nil
}
