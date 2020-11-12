package impl

import (
	"context"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAuthority(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Authority{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	data.ID = NewUUID()
	sc := service.GetContainerInstance()
	err := sc.DB.AddAuthority(context.TODO(), data)
	CheckHandler(err, message.PhoneError)
	HttpReponseHandler(c, data)
}

func GetAuthority(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	data, err := sc.DB.GetAuthority(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, data)
}

func DeleteAuthority(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	err := sc.DB.DeleteAuthority(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, nil)
}

func RoleAuthorityList(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("role_id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	lsit, err := sc.DB.LoadRoleToAuthorityRole(context.TODO(), oid)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, lsit)
}

func AddRole(c *framework.Context) {
	defer PanicHandler(c)
	role := &model.Role{}
	CheckHandler(!c.ParseBody(role), message.JsonParseError)
	sc := service.GetContainerInstance()
	role.ID = NewUUID()
	err := sc.DB.AddRole(context.TODO(), role)
	CheckHandler(err, message.PhoneError)
	HttpReponseHandler(c, role)
}

func DeleteRole(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	err = sc.DB.DeleteRole(ctx, oid)
	CheckHandler(err, message.DeleteError)
	err = sc.DB.DeleteRoleToAuthorityRole(ctx, oid)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func GetRole(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	role, err := sc.DB.GetRole(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, role)
}

func UserRoleList(c *framework.Context) {
	defer PanicHandler(c)
	userID := c.GetParam("user_id")
	userOID, _ := primitive.ObjectIDFromHex(userID)
	institutionID := c.GetParam("institution_id")
	institutionOID, _ := primitive.ObjectIDFromHex(institutionID)
	sc := service.GetContainerInstance()
	lsit, err := sc.DB.LoadUserToInstitutionToRoleUserInstitution(context.TODO(), institutionOID, userOID)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, lsit)
}

func InstitutionRoleList(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")
	oid, _ := primitive.ObjectIDFromHex(institutionID)
	option := db.NewOptions()
	option.EQ[db.OptInstitution] = oid
	sc := service.GetContainerInstance()
	lsit, _, err := sc.DB.LoadRole(context.TODO(), option)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, lsit)
}

func EditRole(c *framework.Context) {
	defer PanicHandler(c)
	role := &model.Role{}
	CheckHandler(!c.ParseBody(role), message.JsonParseError)
	sc := service.GetContainerInstance()
	err := sc.DB.UpdateRole(context.TODO(), role)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, role)
}
