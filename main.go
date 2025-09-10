package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


type TemplateRender struct {
	Templates *template.Template
}

func (t *TemplateRender) Render(w io.Writer, name string, data interface{}, c echo.Context) error{ 
	return t.Templates.ExecuteTemplate(w, name, data)
}

func main() {

	
	renderer := &TemplateRender{
		Templates: template.Must(template.ParseGlob("template/*.html")),
	}
	
	e := echo.New()
	e.Renderer = renderer
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Static("/staitc", "static")
	e.GET("/:id", RedirectHandler)
	e.GET("/", IndexHandler)
	e.POST("/submit", SubmitHandler)
	e.DELETE("/:id", DeleteHandler)
	e.Logger.Fatal(e.Start(":8080"))
	
}