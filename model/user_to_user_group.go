package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserToUserGroup struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Manager   bool               `json:"manager" bson:"manager"`
	User      primitive.ObjectID `json:"user" bson:"user"`
	UserName  string             `json:"user_name" bson:"user_name"`
	Group     primitive.ObjectID `json:"group" bson:"group"`
	GroupName string             `json:"group_name" bson:"group_name"`
}
