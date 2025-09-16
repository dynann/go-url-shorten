package model

import (
	"time"
)

type ClickTime struct {
	Date time.Time 				   `json:"date" bson:"date"`
}
type Link struct {
	Id          string `json:"id" bson:"_id,omitempty"`
	Url         string 			   `json:"url" bson:"url"`
	Clicks      int				   `json:"clicks" bson:"clicks"`	
	ClickRecord []ClickTime		   `json:"click_records" bson:"click_record"`
}

type ClickPerHour struct {
	Hour int
	Click int
}
