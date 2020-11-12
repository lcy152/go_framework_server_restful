package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DipperMessage struct {
	ID                    primitive.ObjectID `json:"_id" bson:"_id"`
	Source                string             `json:"source" bson:"source"`
	Sender                primitive.ObjectID `json:"sender" bson:"sender"`
	SenderName            string             `json:"sender_name" bson:"sender_name"`
	SenderInstitution     primitive.ObjectID `json:"sender_institution" bson:"sender_institution"`
	SenderInstitutionName string             `json:"sender_institution_name" bson:"sender_institution_name"`
	RefPatient            primitive.ObjectID `json:"ref_patient" bson:"ref_patient"`
	RefPatientName        string             `json:"ref_patient_name" bson:"ref_patient_name"`
	Receiver              primitive.ObjectID `json:"receiver" bson:"receiver"`
	ReceiverName          string             `json:"receiver_name" bson:"receiver_name"`
	Institution           primitive.ObjectID `json:"institution" bson:"institution"`
	InstitutionName       string             `json:"institution_name" bson:"institution_name"`
	Type                  string             `json:"type" bson:"type"`
	Flag                  string             `json:"flag" bson:"flag"`
	Msg                   string             `json:"msg" bson:"msg"`
	RefOriginData         primitive.ObjectID `json:"ref_origin_data" bson:"ref_origin_data"`
	CreateTime            time.Time          `json:"create_time" bson:"create_time"`
}
