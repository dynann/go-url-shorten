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
	// e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions, http.MethodPatch},
		// AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowHeaders: []string{"*"},
	}))
	e.Static("/staitc", "static")
	e.GET("/", IndexHandler)
	e.GET("/:id", RedirectHandler)
	e.POST("/submit", SubmitHandler)
	e.DELETE("/:id", DeleteHandler)
	e.Logger.Fatal(e.Start(":8080"))
	
}