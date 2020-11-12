package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Institution struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Address     string             `json:"address" bson:"address"`
	Phone       string             `json:"phone" bson:"phone"`
	Photo       string             `json:"photo" bson:"photo"`
	Qrcode      string             `json:"qrcode" bson:"qrcode"`
	Code        string             `json:"code" bson:"code"`
	KeyCode     string             `json:"key_code" bson:"key_code"`
	Type        string             `json:"type" bson:"type"`
	Level       string             `json:"level" bson:"level"`
	Description string             `json:"description" bson:"description"`
	Creator     primitive.ObjectID `json:"creator" bson:"creator"`
	CreateTime  time.Time          `json:"create_time" bson:"create_time"`
}
