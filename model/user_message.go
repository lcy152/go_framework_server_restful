package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserMessage struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	User primitive.ObjectID `json:"user" bson:"user"`
	Data string             `json:"data" bson:"data"`
}
