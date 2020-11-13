package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstitutionApplication struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	Institution       primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName   string             `json:"institution_name" bson:"institution_name"`
	User              primitive.ObjectID `json:"user" bson:"user"`
	UserName          string             `json:"user_name" bson:"user_name"`
	Status            string             `json:"status" bson:"status"`
	UserToInstitution *UserToInstitution `json:"user_to_institution" bson:"user_to_institution"`
	Description       string             `json:"description" bson:"description"`
	CreateTime        time.Time          `json:"create_time" bson:"create_time"`
	OperateTime       time.Time          `json:"operate_time" bson:"operate_time"`
}

const (
	ApplicationStatusWait    = "wait"
	ApplicationStatusApprove = "approve"
	ApplicationStatusReject  = "reject"
)
