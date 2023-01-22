package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h *Handler) SetAuthCookies(c *gin.Context, userId int) error {
	// add cookie extension for security
	c.SetSameSite(http.SameSite(2))

	// set new access cookie TODO: add settings to Handler: viper must not be here
	accessLifetime := viper.GetInt("cookies.access_lifetime")
	accessToken, err := h.services.Authorization.GenerateToken(
		time.Now().Add(time.Duration(accessLifetime)*time.Second),
		userId,
	)
	if err != nil {
		return fmt.Errorf("can not generate access token: %w", err)
	}

	c.SetCookie(
		viper.GetString("cookies.access_name"),
		accessToken,
		accessLifetime,
		"/",
		viper.GetString("domain"),
		false, // true for https
		true,
	)

	// set new refresh cookie
	refreshLifetime := viper.GetInt("cookies.refresh_lifetime")
	refreshToken, err := h.services.Authorization.GenerateToken(
		time.Now().Add(time.Duration(refreshLifetime)*time.Second),
		userId,
	)
	if err != nil {
		return fmt.Errorf("can not generate refresh token: %w", err)
	}

	c.SetCookie(
		viper.GetString("cookies.refresh_name"),
		refreshToken,
		refreshLifetime,
		"/",
		// TODO: viper must not be here
		viper.GetString("domain"),
		false, // true for https
		true,
	)

	return nil
}

func (h *Handler) ClearAuthCookies(c *gin.Context) error {
	// clear access cookie
	c.SetCookie(
		viper.GetString("cookies.access_name"),
		"accessToken",
		-1,
		"/",
		viper.GetString("domain"),
		false, // true for https
		true,
	)

	// clear refresh cookie
	c.SetCookie(
		viper.GetString("cookies.refresh_name"),
		"refreshToken",
		-1,
		"/",
		viper.GetString("domain"),
		false, // true for https
		true,
	)

	return nil
}
