package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserToInstitution struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Institution     primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName string             `json:"institution_name" bson:"institution_name"`
	Current         bool               `json:"current" bson:"current"`
	User            primitive.ObjectID `json:"user" bson:"user"`
	UserName        string             `json:"user_name" bson:"user_name"`
	Flag            string             `json:"flag" bson:"flag"`
	Type            string             `json:"type" bson:"type"`
	Job             string             `json:"job" bson:"job"`
	DipperUser      string             `json:"dipper_user" bson:"dipper_user"`
	DipperPassword  string             `json:"dipper_password" bson:"dipper_password"`
}
