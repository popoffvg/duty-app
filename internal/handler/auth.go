package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitalygudza/duty-app/internal/model"
)

// @Summary Sign-up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		// todo: create custom errors.New() in repository and check here (for status badrequest if username is exist)
		// todo: swagger description
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.SetAuthCookies(c, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"user id (registered)": userId})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body model.SignInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorMessage
// @Failure 500 {object} errorMessage
// @Failure default {object} errorMessage
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var ok bool
	var input model.SignInInput

	input.Username, input.Password, ok = c.Request.BasicAuth()
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "No basic auth present")
		return
	}

	user, err := h.services.Authorization.CheckCredentials(input)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.SetAuthCookies(c, user.Id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"user id (authorized)": user.Id})
}

func (h *Handler) logout(c *gin.Context) {
	h.ClearAuthCookies(c)
	c.JSON(http.StatusOK, statusMessage{
		Status: "logout successful",
	})
}
