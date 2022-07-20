package router

import (
	"github.com/labstack/echo"

	"zoom-api/controller"
)

func Handler() *echo.Echo {
	// echoのインスタンスを作成
	e := echo.New()
	// ルートを設定
	e.GET("/user/:id", controller.ShowUser())
	e.POST("/user", controller.CreateUser())
	return e
}
