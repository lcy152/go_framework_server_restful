package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	MQUserValidation = "MQUserValidation"
	MQTask           = "MQTask"
	MQMessage        = "MQMessage"
)

type MqMessage struct {
	Flag    string `json:"flag" bson:"flag"`
	KeyCode string `json:"key_code"`
	Data    string `json:"data" bson:"data"`
}

type ValidateDipperUser struct {
	Institution primitive.ObjectID `json:"institution" bson:"institution"`
	Flag        string             `json:"flag" bson:"flag"`
	UserGuid    string             `json:"user_guid" bson:"user_guid"`
	DipperUser  string             `json:"dipper_user" bson:"dipper_user"`
	Password    string             `json:"password" bson:"password"`
	Pass        bool               `json:"pass" bson:"pass"`
}
