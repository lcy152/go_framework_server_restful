package model

type AppConfig struct {
	Guid       string `json:"guid" bson:"_id"`
	AppVersion string `json:"app_version" bson:"app_version"`
	AppUrl     string `json:"app_url" bson:"app_url"`
}
