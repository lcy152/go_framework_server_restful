package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserToUser struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	User       primitive.ObjectID `json:"user" bson:"user"`
	UserName   string             `json:"user_name" bson:"user_name"`
	Friend     primitive.ObjectID `json:"friend" bson:"friend"`
	FriendName string             `json:"friend_name" bson:"friend_name"`
}
