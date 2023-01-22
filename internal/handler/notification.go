package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Router /api/teams/:id/test-notify [post]
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

	team, err := h.services.Team.(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, team)
}
