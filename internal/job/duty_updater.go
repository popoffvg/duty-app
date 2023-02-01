package job

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitalygudza/duty-app/internal/model"
	"github.com/vitalygudza/duty-app/internal/service"
)

type DutyUpdater struct {
	services  *service.Service
	timeNow   time.Time
	debugMode bool
}

func NewDutyUpdater(services *service.Service, debugMode bool) *DutyUpdater {
	return &DutyUpdater{
		services:  services,
		timeNow:   time.Now(),
		debugMode: debugMode,
	}
}

func (du *DutyUpdater) Run() {
	timeNow := time.Now()
	if du.debugMode {
		timeNow = du.timeNow
		// update day for next mock iteration
		du.timeNow = du.timeNow.AddDate(0, 0, 1)
	}

	logrus.Infof("run duty updater for date: %s", timeNow.Format("02-01-2006"))

	newDuties, err := du.services.Duty.UpdateDutiesJob(timeNow)
	if err != nil {
		logrus.Warnf("duty updater job failed: %s", err.Error())
	}

	if newDuties == nil {
		logrus.Errorf("duty updater job failed: new duties is empty")
	}

	du.sendNotifications(newDuties)

	logrus.Info("duties updated")
}

func (du *DutyUpdater) sendNotifications(newDuties []model.Notification) {
	var msg string
	for _, duty := range newDuties {
		if duty.IsDaily {
			msg = fmt.Sprintf("Cегодня дежурит @%s", duty.TeammateName)
		} else {
			msg = fmt.Sprintf("На этой неделе on-call дежурит @%s", duty.TeammateName)
		}

		err := du.services.SendNotification(duty.TeamChannelName, msg)
		if err != nil {
			logrus.Warnf("can not send notification for Team %q: %q, because error occured: %s", duty.TeamName, msg, err.Error())
		}
	}
}
