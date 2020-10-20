package main

import (
	proc "dot-alerts/alerts-processor"
	r "dot-alerts/api"
	"dot-alerts/config"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// setup env
	config.EnvSetup()

	// MQTT Client
	go proc.MQTTClient()

	e := echo.New()

	r.InitRoutes(e)

	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(os.Getenv("ECHO_IP") + ":" + os.Getenv("ECHO_PORT")))

}
