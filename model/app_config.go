package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppConfig struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	AppVersion string             `json:"app_version" bson:"app_version"`
	AppURL     string             `json:"app_url" bson:"app_url"`
}
