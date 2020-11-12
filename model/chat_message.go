package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Sender       primitive.ObjectID `json:"sender" bson:"sender"`
	SenderName   string             `json:"sender_name" bson:"sender_name"`
	Receiver     primitive.ObjectID `json:"receiver" bson:"receiver"`
	ReceiverName string             `json:"receiver_name" bson:"receiver_name"`
	Group        primitive.ObjectID `json:"group" bson:"group"`
	GroupName    string             `json:"group_name" bson:"group_name"`
	Type         string             `json:"type" bson:"type"`
	Flag         string             `json:"flag" bson:"flag"`
	Data         string             `json:"data" bson:"data"`
	CreateTime   time.Time          `json:"create_time" bson:"create_time"`
}

type SingleChatMessage struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Sender       primitive.ObjectID `json:"sender" bson:"sender"`
	SenderName   string             `json:"sender_name" bson:"sender_name"`
	Receiver     primitive.ObjectID `json:"receiver" bson:"receiver"`
	ReceiverName string             `json:"receiver_name" bson:"receiver_name"`
	Type         string             `json:"type" bson:"type"`
	Flag         string             `json:"flag" bson:"flag"`
	Data         string             `json:"data" bson:"data"`
	CreateTime   time.Time          `json:"create_time" bson:"create_time"`
}

type GroupMessage struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Sender     primitive.ObjectID `json:"sender" bson:"sender"`
	SenderName string             `json:"sender_name" bson:"sender_name"`
	Group      primitive.ObjectID `json:"group" bson:"group"`
	GroupName  string             `json:"group_name" bson:"group_name"`
	Type       string             `json:"type" bson:"type"`
	Flag       string             `json:"flag" bson:"flag"`
	Data       string             `json:"data" bson:"data"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
}
