package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserGroup struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Institution     primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName string             `json:"institution_name" bson:"institution_name"`
}
