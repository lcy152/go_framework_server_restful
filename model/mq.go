package model

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
	InstitutionId string `json:"institution_id" bson:"institution_id"`
	Flag          string `json:"flag" bson:"flag"`
	UserGuid      string `json:"user_guid" bson:"user_guid"`
	DipperUser    string `json:"dipper_user" bson:"dipper_user"`
	Password      string `json:"password" bson:"password"`
	Pass          bool   `json:"pass" bson:"pass"`
}
