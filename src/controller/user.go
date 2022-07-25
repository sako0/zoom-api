package controller

import (
	"crypto/tls"
	"net/http"
	"zoom-api/model"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

//ユーザを作成する
func User() echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := session.Get("session", c)
		profile := session.Values["profile"]
		return c.JSON(http.StatusOK, profile)
	}
}

//ユーザをgetする
func FindOrCreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// ssl認証を無視する
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		user := new(model.User)
		if err := c.Bind(user); err != nil {
			return err
		}
		user.UserCreate()

		return c.JSON(http.StatusOK, user)
	}
}
