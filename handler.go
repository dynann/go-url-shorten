package main

import (
	"fmt"
	"math/rand"
	"net/http"

	// "strconv"
	"time"

	"github.com/labstack/echo/v4"
)


type ClickTime struct {
	Date time.Time
}

type Response struct {
	Message string
}

type Link struct {
	Id string
	Url string
	Clicks int
	ClickRecord []ClickTime
}

type Links []Link

var linkMap = map[string]*Link{ "example": { Id: "example", Url: "https://example.com", Clicks: 0, ClickRecord: nil }, }

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"


func RedirectHandler(c echo.Context) error { 
	id := c.Param("id")
	fmt.Println("--> ✅✅✅✅",id)
	link, found := linkMap[id]
	if !found { 
		return  c.String(http.StatusNotFound, "link not found")
	}
	if onClickCheck(link.Url, c) != nil {
		return c.String(http.StatusNotFound, "link is not found")
	}

	// fmt.Println("--> ✅✅✅✅", *link)
	if link.ClickRecord == nil {
		link.ClickRecord = []ClickTime{}
	}

	link.ClickRecord = append(link.ClickRecord, ClickTime{
		Date: time.Now(),
	})
	link.Clicks += 1
	return c.JSON(http.StatusAccepted, *link)
	// return c.Redirect(http.StatusMovedPermanently, link.Url)
}

// func DirectHandler(c echo.Context) error {
// 	id := c.Param("id")
// 	link, found := linkMap[id]
// 	if !found {
// 		return c.String(http)
// 	}
// }


func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []byte
	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	newResult := string(result)
	return newResult
}

//check the link availability before shortening

func SubmitHandler(c echo.Context) error {
	// url := c.FormValue("url")
	req := new(Link)
	if url1 := c.Bind(req); url1 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}
	// if url == "" {
	// 	return c.String(http.StatusBadRequest, "url is required")
	// }

	if !linkValidation(req.Url) {
		return c.String(http.StatusGatewayTimeout, "link is not availabel")
	}

	if !(len(req.Url) >= 4 && (req.Url[:4] == "http"  || req.Url[:5] == "https")) {
		req.Url = "https://" + req.Url
	}

	id := generateRandomString(8)
	// fmt.Println("❤️❤️==> id:", id)
	linkMap[id] = &Link{
		Id:  id,
		Url: req.Url,
	}
	// fmt.Println("hello world⭐😉😂🤣🫤💀✅😚🥀😼😁🗣️", req.Url, req.Id)

	return c.JSON(http.StatusOK,linkMap[id])
}

func IndexHandler(c echo.Context) error {
	links := Links{}

	for _, val := range linkMap {
		links = append(links, *val)
	}
	return c.JSON(http.StatusOK, links)

}


func onClickCheck(url string, c echo.Context) error{
	err := c.Redirect(http.StatusFound, url)
	return err
}


func DeleteHandler(c echo.Context) error {
	req := c.Param("id")
	_, found := linkMap[req]
	if(!found) {
		return c.String(http.StatusNotFound,"link not found")
	}
	newLinkMap := map[string]*Link{}
	for _, link := range linkMap {
		if link.Id != req {
			newLinkMap[link.Id] = &Link{ Id: link.Id, Url: link.Url, Clicks: link.Clicks, ClickRecord: link.ClickRecord}
		}
	}
	linkMap = newLinkMap
	links := Links{}

	for _, val := range linkMap {
		links = append(links, *val)
	}
	return c.JSON(http.StatusOK, links)
}