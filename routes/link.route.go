package routes

import (
	"github.com/dynann/url-shorten/controller"
	"github.com/labstack/echo/v4"
)

func LinkRoute(e *echo.Echo) {
	e.POST("/links", controller.CreateLink)
	e.GET("/links", controller.GetAllLinks)
	e.GET("/links/:Id", controller.GetLink)
	e.DELETE("/links/:Id", controller.DeleteLink)
	e.GET("/links/redirect/:Id", controller.RequestReDirectLink)
}