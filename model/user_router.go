package model

import "time"

type UserRouter struct {
	Guid            string    `json:"guid" bson:"_id"`
	InstitutionId   string    `json:"institution_id" bson:"institution_id"`
	InstitutionName string    `json:"institution_name" bson:"institution_name"`
	IsCurrent       bool      `json:"is_current" bson:"is_current"`
	UserGuid        string    `json:"user_guid" bson:"user_guid"`
	Flag            string    `json:"flag" bson:"flag"`
	Type            string    `json:"type" bson:"type"`
	Job             string    `json:"job" bson:"job"`
	RoleIdList      []string  `json:"role_id_list" bson:"role_id_list"`
	DipperUser      string    `json:"dipper_user" bson:"dipper_user"`
	DipperPassword  string    `json:"dipper_password" bson:"dipper_password"`
	Creator         string    `json:"creator" bson:"creator"`
	LastOperator    string    `json:"last_operator" bson:"last_operator"`
	CreatedTime     time.Time `json:"created_time" bson:"created_time"`
	LastModTime     time.Time `json:"last_mod_time" bson:"last_mod_time"`
}
