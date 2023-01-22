package scheduler

import (
	"github.com/sirupsen/logrus"
	"github.com/vitalygudza/duty-app/internal/service"
)

type Notifier interface {
}

type DutyUpdaterJob struct {
	services *service.Service
}

func NewDutyUpdaterJob(services *service.Service) *DutyUpdaterJob {
	return &DutyUpdaterJob{
		services: services,
	}
}

func (dc DutyUpdaterJob) Run() {
	err := dc.services.Duty.UpdateDutiesJob()
	if err != nil {
		logrus.Warnf("duty updater job failed: %s", err.Error())
	}
}
