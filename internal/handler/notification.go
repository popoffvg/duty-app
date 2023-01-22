package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	testMessage = "test message from Duty-app"
)

// @Summary Send test notification in team's Space channel
// @Security ApiKeyAuth
// @Tags notifications
// @Description send test notification in team's Space channel
// @ID test-notification
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/test-notification [post]
func (h *Handler) sendTestNotification(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// if user is not owner, error is occurred. TODO: maybe it needs check - IsOwner(userId, teamId) (bool, error)
	team, err := h.services.Team.Read(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if team.SpaceChannel != "" {
		err = h.services.Notifier.SendNotification(team.SpaceChannel, testMessage)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("test notification was send to channel: %s", team.SpaceChannel))
	}

	errorResponse(c, http.StatusInternalServerError, "can't send notification: team channel is empty")
}
