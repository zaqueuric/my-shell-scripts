package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {

	//Main route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Dot Alerts Processing")
	})

	//routes.SetEntityRoutes(e)

}
