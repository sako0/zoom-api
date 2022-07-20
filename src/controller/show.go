package controller

import (
	"net/http"
	"zoom-api/model"

	"github.com/labstack/echo"
)

// ユーザを表示する
func ShowUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := model.User{}
		id := c.Param("id")
		targetUser, err := user.FindUserById(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, targetUser)
	}
}

//ユーザを作成する
func CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(model.User)
		if err := c.Bind(user); err != nil {
			return err
		}
		err := user.UserCreate()
		if err != nil {
			return c.JSON(http.StatusOK, err)

		}
		return c.JSON(http.StatusOK, user)

	}
}
