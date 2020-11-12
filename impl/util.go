package impl

import (
	framework "tumor_server/framework"
	service "tumor_server/service"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUUID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func NewUUIDStr() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

func GetContextUserInfo(c *framework.Context) *service.UserTokenInfo {
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	return userInfo
}
