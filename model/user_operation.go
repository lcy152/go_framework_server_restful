package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOperation struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Institution     primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName string             `json:"institution_name" bson:"institution_name"`
	User            primitive.ObjectID `json:"user" bson:"user"`
	UserName        string             `json:"user_name" bson:"user_name"`
	Type            string             `json:"type" bson:"type"`
	RefData         primitive.ObjectID `json:"ref_data" bson:"ref_data"`
	CreateTime      time.Time          `json:"create_time" bson:"create_time"`
	IP              string             `json:"ip" bson:"ip"`
}

const (
	UserOperationLogin       = "login"
	UserOperationLogout      = "logout"
	UserOperationOtherLogin  = "other_login"
	UserOperationOtherLogout = "other_logout"
)
