package model

import "time"

type Notification struct {
	TeamName        string
	TeamChannelName string
	TeammateName    string
	IsDaily         bool
	Date            time.Time
}
