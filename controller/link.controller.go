package controller

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dynann/url-shorten/lib"
	"github.com/dynann/url-shorten/model"
	"github.com/dynann/url-shorten/res"
	"github.com/dynann/url-shorten/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()
const charset = "abcdefghijklmnopqrstuvwxyz0123456789"


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

func onClickCheck(url string, c echo.Context) error{
	err := c.Redirect(http.StatusFound, url)
	return err
}

func countClickPerHour(clickRecord []model.ClickTime) []model.ClickPerHour{
	today := time.Now().Day()
	if len(clickRecord) == 0 {
		return []model.ClickPerHour{}
	}
	values := []model.ClickPerHour{}
	hour := 0
	for hour < 24 {
		var counts int = 0
		for _, clicks := range clickRecord {
			if(clicks.Date.Day() == today) {
				if (clicks.Date.Hour() + 7) == hour {
				counts += 1
			}
			}	
		}
		if (hour == 0) {
			data := model.ClickPerHour{
			Hour: 12,
			Click: counts,
		} 
		values = append(values, data)
		} else if (hour > 12) {
			data := model.ClickPerHour{
				Hour: (hour - 12),
				Click: counts,
			}
		values = append(values, data)
		} else if (hour < 12) {
			data := model.ClickPerHour{
				Hour: hour,
				Click: counts,
			}
		values = append(values, data)
		} else {
			data := model.ClickPerHour{
				Hour: hour,
				Click: counts,
			}
		values = append(values, data)
		}
		hour += 1
	}
	return values
}
 

func GetLink(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	linkId := c.Param("Id")
	fmt.Println(linkId)
	var link model.Link
	defer cancel()

	err := linkCollection.FindOne(ctx, bson.M{"_id": linkId}).Decode(&link)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, res.LinkResponse{
	    Status: http.StatusOK,
		Message: "success",
		Data: &echo.Map{"data": link}})
}

func GetAllLinks(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var links []model.Link
	defer cancel()

	values, err := linkCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{"data": err.Error()},
		})
	}

	if err := values.All(ctx, &links); err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{"data": err.Error()},
		})
	}

	if len(links) == 0 {
		return c.JSON(http.StatusOK, res.LinkResponse{
		Status: http.StatusAccepted,
		Message: "error",
		Data: &echo.Map{"data": []model.Link{} },
	})
	}


	return c.JSON(http.StatusOK, res.LinkResponse{
		Status: http.StatusAccepted,
		Message: "error",
		Data: &echo.Map{"data": links },
	})
}

func CreateLink(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var link model.Link
	defer cancel()

	if err := c.Bind(&link); err != nil {
		return c.JSON(http.StatusBadRequest, res.LinkResponse{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&link); validationErr != nil {
		return c.JSON(http.StatusBadRequest, res.LinkResponse{
			Status: http.StatusBadRequest,
			Message: "Error",
			Data: &echo.Map{"data": validationErr.Error()}})
	}

	if !validation.LinkValidation(link.Url) {
		return c.JSON(http.StatusRequestTimeout, res.LinkResponse{
			Status: http.StatusRequestTimeout, 
			Message: "error time out", 
			Data: &echo.Map{"data": nil}})
	}

	if !(len(link.Url) >= 4 && (link.Url[:4] == "http" || link.Url[:5] == "https")) {
		link.Url = "https://" + link.Url
	}

	link.Id = generateRandomString(8)

	newLink := model.Link{
		Id:          link.Id,
		Url:         link.Url,
		Clicks:      link.Clicks,
		ClickRecord: []model.ClickTime{},
	}

	result, err := linkCollection.InsertOne(ctx, newLink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
	}

	fmt.Print("-->>>>😁😁😼😼", result)

	return c.JSON(http.StatusCreated, res.LinkResponse{
		Status: http.StatusCreated, 
		Message: "success", 
		Data: &echo.Map{"data": newLink}})

}

func DeleteLink(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	linkId := c.Param("Id")
	result, err := linkCollection.DeleteOne(ctx, bson.M{"_id": linkId})
	fmt.Println("❤️❤️❤️", err, "❤️❤️❤️")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError,
			Message: "error",
			Data: &echo.Map{"data": err.Error()},
		})
	}
	return c.JSON(http.StatusOK, res.LinkResponse{
		Status: http.StatusOK,
		Message: "successfully deleted link",
		Data: &echo.Map{"data": result},

	} )
}

func RequestReDirectLink(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var link model.Link

	defer cancel()

	linkId := c.Param("Id")
	err := linkCollection.FindOne(ctx, bson.M{"_id": linkId}).Decode(&link)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError,
			Message: "internal server error",
			Data: &echo.Map{"data": err.Error()},
		})
	}

	if onClickCheck(link.Url, c) != nil {
		return c.String(http.StatusNotFound, "link is not found")
	}

	if link.ClickRecord == nil {
		link.ClickRecord = []model.ClickTime{}
	}



	link.ClickRecord = append(link.ClickRecord, model.ClickTime{
		Date: time.Now(),
	})
	link.Clicks += 1
	err = updateOnLinkClick(link)
	fmt.Println("-->🙃😂🫠", err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError,
			Message: "fail to update",
			Data: &echo.Map{"data": err.Error()},
		})
	}
	return c.JSON(http.StatusAccepted, res.LinkResponse{
			Status: http.StatusAccepted,
			Message: "success",
			Data: &echo.Map{"data": link},
		})

}


func updateOnLinkClick(link model.Link) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := linkCollection.FindOneAndUpdate(ctx, bson.M{"_id": link.Id}, bson.M{
		"$inc": bson.M{"clicks": 1},
		"$push": bson.M{"click_record": bson.M{
			"date": time.Now(),
		}},
	}, opts)

	if err != nil {
		return err.Err()
	}

	return nil
}


func ReturnClickByHours(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Param("Id")
	var link model.Link
	defer cancel()
	err := linkCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&link)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
	}
	var data = countClickPerHour(link.ClickRecord)
	// fmt.Println("==> ❤️❤️", data)
	return c.JSON(http.StatusOK, res.LinkResponse{
			Status: http.StatusOK, 
			Message: "success", 
			Data: &echo.Map{"data": data}})
}

func RequestDirectLink(c echo.Context) error {
	linkCollection := lib.GetCollection("links")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Param("Id")
	link := model.Link{}
	defer cancel()
	
	err := linkCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&link)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.LinkResponse{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &echo.Map{"data": err.Error()}})
	}
	if onClickCheck(link.Url, c) != nil {
		return c.String(http.StatusNotFound, "link is not found")
	}

	if link.ClickRecord == nil {
		link.ClickRecord = []model.ClickTime{}
	}



	link.ClickRecord = append(link.ClickRecord, model.ClickTime{
		Date: time.Now(),
	})
	link.Clicks += 1
	err = updateOnLinkClick(link)
	if err != nil {
	return c.JSON(http.StatusInternalServerError, res.LinkResponse{
		Status: http.StatusInternalServerError,
		Message: "fail to update",
		Data: &echo.Map{"data": err.Error()},
	})
	}
	return c.JSON(http.StatusAccepted, res.LinkResponse{
		Status: http.StatusAccepted,
		Message: "successful",
		Data: &echo.Map{"data":  link},
	})
}