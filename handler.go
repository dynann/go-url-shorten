package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)


type ClickTime struct {
	Date time.Time
}

type Link struct {
	Id string
	Url string
	Clicks int
	ClickRecord []ClickTime
}

var linkMap = map[string]*Link{ "example": { Id: "example", Url: "https://example.com", Clicks: 0, ClickRecord: nil }, }

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
func RedirectHandler(c echo.Context) error { 
	id := c.Param("id")
	link, found := linkMap[id]
	if !found { 
		return  c.String(http.StatusNotFound, "link not found")
	}
	if onClickCheck(link.Url, c) != nil {
		return c.String(http.StatusNotFound, "link not found")
	}
	if link.ClickRecord == nil {
		link.ClickRecord = []ClickTime{}
	}
	link.ClickRecord = append(link.ClickRecord, ClickTime{
		Date: time.Now(),
	})
	link.Clicks += 1
	fmt.Println("here is the value ✅✅", link.ClickRecord)
	return c.Redirect(http.StatusMovedPermanently, link.Url)
}


func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []byte
	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	newResult := "https://" + string(result) + ".dn"
	return newResult
}

//check the link availability before shortening

func SubmitHandler(c echo.Context) error {
	url := c.FormValue("url")
	if url == "" {
		return c.String(http.StatusBadRequest, "url is required")
	}

	if !linkValidation(url) {
		return c.String(http.StatusGatewayTimeout, "link is not availabel")
	}

	if !(len(url) >= 4 && (url[:4] == "http"  || url[:5] == "https")) {
		url = "https://" + url
	}

	id := generateRandomString(8)
	linkMap[id] = &Link{
		Id: id,
		Url: url,
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func IndexHandler(c echo.Context) error {
	html := `
		<h1>Submit a new website</h1>
		<form action="/submit" method="POST">
		<label for="url">Website URL:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Submit">
		</form>
		<h2>Existing Links </h2>
		<ul>`

	for _, link := range linkMap {
		// Show the link and click count
		html += `<li><a href="/` + link.Id + `">` + link.Id + `</a> ` + strconv.Itoa(link.Clicks) + `</li>`

		// Start table
		html += `<table border="1" cellpadding="5">
		<tr><th>Click Time</th></tr>`

		// Loop through each click record
		for _, record := range link.ClickRecord {
			html += `<tr><td>` + record.Date.Format("2006-01-02 15:04:05") + `</td></tr>`
		}

		html += `</table>`
	}

	html += `</ul>`
	return c.HTML(http.StatusOK, html)
}


func onClickCheck(url string, c echo.Context) error{
	err := c.Redirect(http.StatusFound, url)
	return err
}