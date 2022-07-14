package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vitalygudza/duty-app/internal/model"
)

// @Summary Create team
// @Security ApiKeyAuth
// @Tags teams
// @Description create new team
// @ID create-team
// @Accept  json
// @Produce  json
// @Param input body model.Team true "team info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams [post]
func (h *Handler) createTeam(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.Team
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Team.Create(userId, input)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// @Summary Read team
// @Security ApiKeyAuth
// @Tags teams
// @Description Read team by id
// @ID read-team
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Team
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id [get]
func (h *Handler) readTeam(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	team, err := h.services.Team.Read(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, team)
}

// @Summary Update team
// @Security ApiKeyAuth
// @Tags teams
// @Description Update team by id
// @ID update-team
// @Accept  json
// @Produce  json
// @Param input body model.UpdateTeamInput true "team info"
// @Success 200 {string} string "team updated"
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id [put]
func (h *Handler) updateTeam(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.UpdateTeamInput
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Team.Update(userId, teamId, input); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{
		Status: "team updated",
	})
}

// @Summary Delete team
// @Security ApiKeyAuth
// @Tags teams
// @Description Delete team by id
// @ID delete-team
// @Accept  json
// @Produce  json
// @Success 200 {string} string "team deleted"
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id [delete]
func (h *Handler) deleteTeam(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Team.Delete(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{
		Status: "team deleted",
	})
}

// @Summary List teams
// @Security ApiKeyAuth
// @Tags teams
// @Description List all teams
// @ID list-teams
// @Accept  json
// @Produce  json
// @Success 200 {object} listTeamsResponse
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams [get]
func (h *Handler) listTeams(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teams, err := h.services.Team.List(userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listTeamsResponse{
		Data: teams,
	})
}
