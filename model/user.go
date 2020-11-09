package model

import "time"

const (
	AdminGuid = "admin"
)

type Contacter struct {
	Name         string `json:"name" bson:"name"`
	Phone        string `json:"phone" bson:"phone"`
	FixedPhone   string `json:"fixed_phone" bson:"fixed_phone"`
	Relationship string `json:"relationship" bson:"relationship"`
}

type User struct {
	Guid        string        `json:"guid" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Phone       string        `json:"phone" bson:"phone"`
	Password    string        `json:"password" bson:"password"`
	Token       string        `json:"token" bson:"token"`
	Disable     bool          `json:"disable" bson:"disable"`
	Status      string        `json:"status" bson:"status"`
	Type        string        `json:"type" bson:"type"`
	Sex         string        `json:"sex" bson:"sex"`
	BirthDate   string        `json:"birth_date" bson:"birth_date"`
	IDCard      string        `json:"id_card" bson:"id_card"`
	Address     string        `json:"address" bson:"address"`
	Height      string        `json:"height" bson:"height"`
	Weight      string        `json:"weight" bson:"weight"`
	Qrcode      string        `json:"qrcode" bson:"qrcode"`
	Photo       string        `json:"photo" bson:"photo"`
	Contacter   []*Contacter  `json:"contacter" bson:"contacter"`
	RouterList  []*UserRouter `json:"router_list" bson:"router_list"`
	FriendList  []string      `json:"friend_list" bson:"friend_list"`
	CreateTime  time.Time     `json:"create_time" bson:"create_time"`
	LastModTime time.Time     `json:"last_mod_time" bson:"last_mod_time"`
}

type Guid struct {
	Guid string `json:"guid"`
}
