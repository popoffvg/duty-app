package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	var userId int
	var err error

	userId, err = h.getUserIdFromCookie(c, viper.GetString("cookies.access_name"))
	if err != nil {
		// if access cookie is expired, try to get userId from refresh cookie and refresh the cookies
		userId, err = h.getUserIdFromCookie(c, viper.GetString("cookies.refresh_name"))
		if err != nil {
			errorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		// update access and refresh cookie
		err = h.SetAuthCookies(c, userId)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Set(userCtx, userId)
	c.Next()
}

func (h *Handler) getUserIdFromCookie(c *gin.Context, name string) (int, error) {
	token, err := c.Cookie(name)
	if err != nil {
		return 0, fmt.Errorf("can not get cookie: %w", err)
	}

	userId, err := h.services.ParseToken(token)
	if err != nil {
		return 0, fmt.Errorf("can not parse %q cookie: %w", name, err)
	}

	return userId, nil
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, "user id is invalid type")
		return 0, errors.New("user id is invalid type")
	}

	return idInt, nil
}
