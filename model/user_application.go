package model

import "time"

type UserApplication struct {
	Guid            string    `json:"guid" bson:"_id"`
	InstitutionId   string    `json:"institution_id" bson:"institution_id"`
	InstitutionName string    `json:"institution_name" bson:"institution_name"`
	UserGuid        string    `json:"user_guid" bson:"user_guid"`
	Type            string    `json:"type" bson:"type"`
	Status          string    `json:"status" bson:"status"`
	Creator         string    `json:"creator" bson:"creator"`
	CreateTime      time.Time `json:"create_time" bson:"create_time"`
	Description     string    `json:"description" bson:"description"`
}

const (
	ApplicationTypeApplyFriend      = "apply_friend"
	ApplicationTypeApplyInstitution = "apply_institution"
)

const (
	ApplicationStatusWait    = "wait"
	ApplicationStatusApprove = "approve"
	ApplicationStatusReject  = "reject"
)
