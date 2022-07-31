package router

import (
	"crypto/tls"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"zoom-api/controller"
)

func Handler() *echo.Echo {
	// echoのインスタンスを作成
	e := echo.New()

	// CORSの設定追加
	e.Use(middleware.CORS())

	// ssl認証を無視する
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// ルートを設定
	e.POST("/user", controller.FindOrCreateUser())
	e.POST("/login", controller.Login())
	e.POST("/createZoom", controller.CreateZoom())
	e.GET("/getMyZoomList", controller.MyZoomList())

	return e
}
