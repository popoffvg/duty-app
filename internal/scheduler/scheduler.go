package cron

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
}

func (s *Scheduler) NewScheduler() {
	var scheduler *cron.Cron

	locationTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return
	}

	scheduler = cron.New(cron.WithLocation(locationTime))
	scheduler.AddFunc("", s.UpdateDuties)

	go scheduler.Run()
}

func (s *Scheduler) UpdateDuties() {

}
