package impl

import (
	"context"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	ug := model.UserGroup{}
	CheckHandler(!c.ParseBody(&ug), message.JsonParseError)
	ug.ID = NewUUID()
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.UnknownError)
	defer session.Close()
	err = sc.DB.AddUserGroup(ctx, &ug)
	CheckHandler(err, message.AddError)
	userInfo := GetContextUserInfo(c)
	{
		utu := &model.UserToUserGroup{}
		utu.ID = NewUUID()
		utu.Group = ug.ID
		utu.GroupName = ug.Name
		utu.User = userInfo.User.ID
		utu.UserName = userInfo.User.Name
		utu.Manager = true
		err := sc.DB.AddUserToUserGroup(ctx, utu)
		CheckHandler(err, message.AddError)
	}
	session.Commit()
	HttpReponseHandler(c, ug)
}

func GetUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	group, err := sc.DB.GetUserGroup(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, group)
}

func DeleteUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	err := sc.DB.DeleteUserGroup(context.TODO(), oid)
	CheckHandler(err, message.DeleteError)
	err = sc.DB.DeleteUserToUserGroupGroup(context.TODO(), oid)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func EditUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.UserGroup{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	err := sc.DB.UpdateUserGroup(context.TODO(), data)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, data)
}

func UserUserGroupList(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	ugList, err := sc.DB.LoadUserToUserGroupUser(context.TODO(), userInfo.User.ID)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, ugList)
}

func AddUserGroupUser(c *framework.Context) {
	defer PanicHandler(c)
	data := []*model.UserToUserGroup{}
	CheckHandler(!c.ParseBody(&data), message.JsonParseError)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.UnknownError)
	defer session.Close()
	for _, v := range data {
		v.ID = NewUUID()
		err := sc.DB.AddUserToUserGroup(ctx, v)
		CheckHandler(err, message.UpdateError)
	}
	session.Commit()
	HttpReponseHandler(c, data)
}

func GetUserGroupUser(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	utuList, err := sc.DB.LoadUserToUserGroupGroup(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, utuList)
}

func DeleteUserGroupUser(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("group_id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	err := sc.DB.DeleteUserToUserGroup(context.TODO(), oid)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, nil)
}
