package model

import "time"

type UserGroup struct {
	Guid       string    `json:"guid" bson:"_id"`
	Name       string    `json:"name" bson:"name"`
	Manager    []string  `json:"manager" bson:"manager"`
	Member     []string  `json:"member" bson:"member"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}
