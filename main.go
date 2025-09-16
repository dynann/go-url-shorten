package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/dynann/url-shorten/lib"
	"github.com/dynann/url-shorten/routes"
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

	routes.UserRoute(e)
	routes.LinkRoute(e)


	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("error loading .env file")
	// }

	var PORT = os.Getenv("PORT")
	
	if err := lib.InitializeMongoDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)  // ✅ Shows actual error
	}

	// ✅ Connection test
	if !lib.IsConnected() {
		log.Fatal("Database connection test failed")
	}
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions, http.MethodPatch},
		// AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowHeaders: []string{"*"},
	}))
	e.GET("/indirect/:id", RedirectHandler)
	e.Logger.Fatal(e.Start(":"+PORT))
	
}
