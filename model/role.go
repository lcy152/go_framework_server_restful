package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authority struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Number string             `json:"number" bson:"number"`
	Name   string             `json:"name" bson:"name"`
}

type Role struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Institution     primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName string             `json:"institution_name" bson:"institution_name"`
	Name            string             `json:"name" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Status          string             `json:"status" bson:"status"`
	Editable        bool               `json:"editable" bson:"editable"`
	Creator         string             `json:"creator" bson:"creator"`
}

type RoleToAuthority struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Role            primitive.ObjectID `json:"role" bson:"role"`
	RoleName        string             `json:"role_name" bson:"role_name"`
	Authority       primitive.ObjectID `json:"authority" bson:"authority"`
	AuthorityNumber string             `json:"authority_number" bson:"authority_number"`
	AuthorityName   string             `json:"authority_name" bson:"authority_name"`
}
