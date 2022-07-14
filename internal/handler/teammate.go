package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vitalygudza/duty-app/internal/model"
)

// @Summary Create teammate
// @Security ApiKeyAuth
// @Tags teammates
// @Description create new teammate
// @ID create-teammate
// @Accept  json
// @Produce  json
// @Param input body model.Teammate true "teammate info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/teammates [post]
func (h *Handler) createTeammate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.Teammate
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Teammate.Create(userId, teamId, input)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Read teammate
// @Security ApiKeyAuth
// @Tags teammates
// @Description read teammate by id
// @ID read-teammate
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Teammate
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teammates/:id [get]
func (h *Handler) readTeammate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teammateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	teammate, err := h.services.Teammate.Read(userId, teammateId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, teammate)
}

// @Summary Update teammate
// @Security ApiKeyAuth
// @Tags teammates
// @Description update teammate by id
// @ID update-teammate
// @Accept  json
// @Produce  json
// @Param input body model.UpdateTeammateInput true "teammate info"
// @Success 200 {string} string "teammate updated"
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teammates/:id [put]
func (h *Handler) updateTeammate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teammateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.UpdateTeammateInput
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Teammate.Update(userId, teammateId, input); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{
		Status: "teammate updated",
	})
}

// @Summary Delete teammate
// @Security ApiKeyAuth
// @Tags teammates
// @Description delete teammate by id
// @ID delete-teammate
// @Accept  json
// @Produce  json
// @Success 200 {string} string "teammate deleted"
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teammates/:id [delete]
func (h *Handler) deleteTeammate(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teammateId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Teammate.Delete(userId, teammateId); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{Status: "teammate deleted"})
}

// @Summary List teammates
// @Security ApiKeyAuth
// @Tags teammates
// @Description list all teammates
// @ID list-teammates
// @Accept  json
// @Produce  json
// @Success 200 {object} listTeammatesResponse
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/teammates [get]
func (h *Handler) listTeammates(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	teammates, err := h.services.Teammate.List(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listTeammatesResponse{
		Data: teammates,
	})
}
