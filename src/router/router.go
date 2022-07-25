package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"zoom-api/controller"
)

func Handler() *echo.Echo {
	// echoのインスタンスを作成
	e := echo.New()
	// CORSの設定追加
	e.Use(middleware.CORS())

	// ルートを設定
	e.POST("/user", controller.FindOrCreateUser())

	return e
}
