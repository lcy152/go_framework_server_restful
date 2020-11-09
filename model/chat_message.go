package model

import "time"

type ChatMessage struct {
	Guid         string    `json:"guid" bson:"_id"`
	SenderGuid   string    `json:"sender_guid" bson:"sender_guid"`
	SenderName   string    `json:"sender_name" bson:"sender_name"`
	ReceiverGuid string    `json:"receiver_guid" bson:"receiver_guid"`
	ReceiverName string    `json:"receiver_name" bson:"receiver_name"`
	GroupGuid    string    `json:"group_guid" bson:"group_guid"`
	Type         string    `json:"type" bson:"type"`
	Flag         string    `json:"flag" bson:"flag"`
	Data         string    `json:"data" bson:"data"`
	CreateTime   time.Time `json:"create_time" bson:"create_time"`
}
