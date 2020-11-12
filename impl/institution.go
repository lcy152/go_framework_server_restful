package impl

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddInstitution(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Institution{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	tn := time.Now()
	data.ID = NewUUID()
	data.CreateTime = tn
	data.Creator = userInfo.User.ID
	err = sc.DB.AddInstitution(ctx, data)
	CheckHandler(err, message.AddError)
	userRefIns := &model.UserToInstitution{
		ID:              NewUUID(),
		Institution:     data.ID,
		InstitutionName: data.Name,
		User:            userInfo.User.ID,
		UserName:        userInfo.User.Name,
	}
	err = sc.DB.AddUserToInstitution(ctx, userRefIns)
	CheckHandler(err, message.AddError)
	session.Commit()
	if sc.RabbitMQ != nil {
		sc.RabbitMQ.Connect()
	}
	HttpReponseHandler(c, data)
}

func GetInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	ins, err := sc.DB.GetInstitution(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	CheckHandler(ins == nil, message.GetError)
	HttpReponseHandler(c, ins)
}

func EditInstitution(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Institution{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	ins, err := sc.DB.GetInstitution(context.TODO(), data.ID)
	CheckHandler(err, message.GetError)
	CheckHandler(!c.ParseBody(ins), message.JsonParseError)
	err = sc.DB.UpdateInstitution(context.TODO(), ins)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, ins)
}

func DeleteInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.UnknownError)
	defer session.Close()
	err = sc.DB.DeleteInstitution(ctx, oid)
	CheckHandler(err, message.DeleteError)
	err = sc.DB.DeleteUserToInstitutionInstitution(ctx, oid)
	CheckHandler(err, message.DeleteError)
	session.Commit()
	HttpReponseHandler(c, nil)
}

func LoadInstitution(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	search := c.GetURLParam("search")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)

	option := db.NewOptions()
	option.PageSize = pageSize
	option.PageIndex = pageIndex
	if search != "" {
		option.Match[db.OptName] = search
		option.Match[db.OptCode] = search
		option.Match[db.OptAddress] = search
	}
	sc := service.GetContainerInstance()
	insList, count, err := sc.DB.LoadInstitution(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseListHandler(c, count, insList)
}

func ApplyInstitution(c *framework.Context) {
	var data = &model.InstitutionApplication{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	institution, err := sc.DB.GetInstitution(context.TODO(), data.Institution)
	CheckHandler(institution == nil, message.GetError)
	option := db.NewOptions()
	option.EQ[db.OptUser] = userInfo.User.ID
	option.EQ[db.OptInstitution] = institution.ID
	uaList, _, _ := sc.DB.LoadInstitutionApplication(context.TODO(), option)
	for _, v := range uaList {
		CheckHandler(v.Status == model.ApplicationStatusWait, message.RequestRepeatError)
	}
	err = sc.DB.AddInstitutionApplication(context.TODO(), data)
	CheckHandler(err, message.AddError)
	opt := db.NewOptions()
	opt.EQ[db.OptInstitution] = data.Institution
	opt.EQ[db.OptManager] = true
	urList, err := sc.DB.LoadUserToInstitution(context.TODO(), opt)
	for _, v := range urList {
		jsonStr, _ := json.Marshal(data)
		service.SendInstitutionMessage(v.User.String(), string(jsonStr))
	}
	HttpReponseHandler(c, nil)
}

func InstitutionUserList(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")
	oid, _ := primitive.ObjectIDFromHex(institutionID)
	sc := service.GetContainerInstance()
	urList, err := sc.DB.LoadUserToInstitutionInstitution(context.TODO(), oid)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, urList)
}

func ApproveInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	op, err := sc.DB.GetInstitutionApplication(ctx, oid)
	CheckHandler(err, message.GetError)
	err = sc.DB.AddUserToInstitution(ctx, op.UserToInstitution)
	CheckHandler(err, message.AddError)
	op.Status = model.ApplicationStatusApprove
	err = sc.DB.UpdateInstitutionApplication(ctx, op)
	CheckHandler(err, message.UpdateError)
	session.Commit()
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func RejectInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	op, err := sc.DB.GetInstitutionApplication(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	op.Status = model.ApplicationStatusReject
	err = sc.DB.UpdateInstitutionApplication(context.TODO(), op)
	CheckHandler(err, message.UpdateError)
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func InstitutionApply(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")
	oid, _ := primitive.ObjectIDFromHex(institutionID)
	state := c.GetParam("state")
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.EQ[db.OptInstitution] = oid
	if state != "" {
		opt.EQ[db.OptStatus] = state
	}
	opList, _, err := sc.DB.LoadInstitutionApplication(context.TODO(), opt)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, opList)
}

func AddInstitutionUser(c *framework.Context) {
	defer PanicHandler(c)
	data := []*model.UserToInstitution{}
	CheckHandler(!c.ParseBody(&data), message.JsonParseError)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	for _, v := range data {
		v.ID = NewUUID()
		err = sc.DB.AddUserToInstitution(ctx, v)
		CheckHandler(err, message.AddError)
	}
	session.Commit()
	HttpReponseHandler(c, data)
}

func EditInstitutionUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.UserToInstitution{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	err := sc.DB.UpdateUserToInstitution(context.TODO(), data)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, data)
}

func DeleteInstitutionUser(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	err := sc.DB.DeleteUserToInstitution(context.TODO(), oid)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, nil)
}
