package model

import "time"

type DipperMessage struct {
	Guid                  string    `json:"guid" bson:"_id"`
	Source                string    `json:"source" bson:"source"`
	SenderGuid            string    `json:"sender_guid" bson:"sender_guid"`
	SenderName            string    `json:"sender_name" bson:"sender_name"`
	SenderInstitutionId   string    `json:"sender_institution_id" bson:"sender_institution_id"`
	SenderInstitutionName string    `json:"sender_institution_name" bson:"sender_institution_name"`
	RefPatientGuid        string    `json:"ref_patient_guid" bson:"ref_patient_guid"`
	RefPatientName        string    `json:"ref_patient_name" bson:"ref_patient_name"`
	ReceiverGuid          string    `json:"receiver_guid" bson:"receiver_guid"`
	ReceiverName          string    `json:"receiver_name" bson:"receiver_name"`
	InstitutionId         string    `json:"institution_id" bson:"institution_id"`
	InstitutionName       string    `json:"institution_name" bson:"institution_name"`
	Type                  string    `json:"type" bson:"type"`
	Flag                  string    `json:"flag" bson:"flag"`
	Msg                   string    `json:"msg" bson:"msg"`
	RefOriginDataGuid     string    `json:"ref_origin_data_guid" bson:"ref_origin_data_guid"`
	CreateTime            time.Time `json:"create_time" bson:"create_time"`
}
