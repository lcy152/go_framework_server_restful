package model

import "time"

type UserOperation struct {
	Guid       string    `json:"guid" bson:"_id"`
	UserGuid   string    `json:"user_guid" bson:"user_guid"`
	UserName   string    `json:"user_name" bson:"user_name"`
	Type       string    `json:"type" bson:"type"`
	RefGuid    string    `json:"ref_guid" bson:"ref_guid"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	IP         string    `json:"ip" bson:"ip"`
}

const (
	UserOperationLogin       = "login"
	UserOperationLogout      = "logout"
	UserOperationOtherLogin  = "other_login"
	UserOperationOtherLogout = "other_logout"
)
