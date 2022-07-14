package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vitalygudza/duty-app/internal/model"
)

// @Summary Create duty
// @Security ApiKeyAuth
// @Tags duties
// @Description create new duty
// @ID create-duty
// @Accept  json
// @Produce  json
// @Param input body model.Duty true "duty info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/duties [post]
func (h *Handler) createDuty(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.Duty
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	duty, err := h.services.Duty.Create(userId, teamId, input)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var msg string
	if duty.IsDaily {
		msg = fmt.Sprintf("create daily duty successful (%s)", duty.Date.Format("02-01-2006"))
	} else {
		msg = fmt.Sprintf(
			"create weekly duty successful (%s - %s)",
			duty.Date.Format("02-01-2006"),
			duty.Date.AddDate(0, 0, 6).Format("02-01-2006"))
	}

	// todo refactoring responses
	c.JSON(http.StatusOK, statusMessage{
		Id:     duty.Id,
		Status: msg,
	})
}

// @Summary Update duty
// @Security ApiKeyAuth
// @Tags duties
// @Description update duty
// @ID update-duty
// @Accept  json
// @Produce  json
// @Param input body model.UpdateDutyInput true "duty info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/duties/:id [put]
func (h *Handler) updateDuty(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	dutyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input model.UpdateDutyInput
	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Duty.Update(userId, dutyId, input); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{
		Status: "duty updated",
	})
}

// @Summary Delete duty
// @Security ApiKeyAuth
// @Tags duties
// @Description delete duty by id
// @ID delete-duty
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/duties/:id [delete]
func (h *Handler) deleteDuty(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	dutyId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err = h.services.Duty.Delete(userId, dutyId); err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusMessage{Status: "duty deleted"})
}

// @Summary List duties
// @Security ApiKeyAuth
// @Tags duties
// @Description list all duties
// @ID list-duties
// @Accept  json
// @Produce  json
// @Success 200 {object} listDutiesResponse
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/duties [get]
func (h *Handler) listDuties(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	duties, err := h.services.Duty.List(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listDutiesResponse{
		Data: duties,
	})
}

// @Summary Read duties
// @Security ApiKeyAuth
// @Tags duties
// @Description read current team duties
// @ID read-duty
// @Accept  json
// @Produce  json
// @Success 200 {object} listDutiesResponse
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teammates/:id [get]
func (h *Handler) readDuties(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	duties, err := h.services.Duty.ReadCurrent(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listDutiesResponse{
		Data: duties,
	})
}

// @Summary History duties
// @Security ApiKeyAuth
// @Tags duties
// @Description history of top 100 duties in team
// @ID history-duties
// @Accept  json
// @Produce  json
// @Success 200 {object} listDutiesHistoryResponse
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /api/teams/:id/history [get]
func (h *Handler) historyDuties(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	teamId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	history, err := h.services.Duty.History(userId, teamId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, listDutiesHistoryResponse{
		Data: history,
	})
}
