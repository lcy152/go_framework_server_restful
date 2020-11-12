package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserToInstitutionToRole struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Institution primitive.ObjectID `json:"institution" bson:"institution"`
	User        primitive.ObjectID `json:"user" bson:"user"`
	Role        primitive.ObjectID `json:"role" bson:"role"`
}
