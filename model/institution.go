package model

import "time"

type Institution struct {
	Guid          string    `json:"guid" bson:"_id"`
	Name          string    `json:"name" bson:"name"`
	Manager       []string  `json:"manager" bson:"manager"`
	Address       string    `json:"address" bson:"address"`
	Phone         string    `json:"phone" bson:"phone"`
	Photo         string    `json:"photo" bson:"photo"`
	Code          string    `json:"code" bson:"code"`
	KeyCode       string    `json:"key_code" bson:"key_code"`
	Creator       string    `json:"creator" bson:"creator"`
	CreateTime    time.Time `json:"create_time" bson:"create_time"`
}
