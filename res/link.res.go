package res

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClickTimeResponse struct {
	Id   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Date time.Time          `json:"date" bson:"date"`
}
type LinkResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data *echo.Map `json:"data"`
}
