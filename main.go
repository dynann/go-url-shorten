package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/:id", RedirectHandler)
	e.GET("/", IndexHandler)
	e.POST("/submit", SubmitHandler )
	e.Logger.Fatal(e.Start(":8080"))
	
}