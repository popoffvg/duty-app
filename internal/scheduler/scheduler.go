package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Run()
}

type Scheduler struct {
	jobRunner *cron.Cron
}

func NewScheduler() (*Scheduler, error) {
	locationTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("can not set location time in scheduler: %w", err)
	}

	scheduler := Scheduler{
		jobRunner: cron.New(cron.WithLocation(locationTime)),
	}

	return &scheduler, nil
}

func (s *Scheduler) AddJob(schedule string, job Job) error {
	_, err := s.jobRunner.AddJob(schedule, job)
	if err != nil {
		return fmt.Errorf("can not schedule the job: %w", err)
	}

	return nil
}

func (s *Scheduler) Start() {
	s.jobRunner.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.jobRunner.Stop()
}
