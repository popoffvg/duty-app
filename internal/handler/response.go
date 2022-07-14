package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vitalygudza/duty-app/internal/model"
)

type errorMessage struct {
	Message string `json:"message"`
}

type statusMessage struct {
	Id     int    `json:"id,omitempty"`
	Status string `json:"status"`
}

type listTeamsResponse struct {
	Data []model.Team `json:"data"`
}

type listTeammatesResponse struct {
	Data []model.Teammate `json:"data"`
}

type listDutiesResponse struct {
	Data []model.Duty `json:"data"`
}

type listDutiesHistoryResponse struct {
	Data []model.History `json:"data"`
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorMessage{message})
}
