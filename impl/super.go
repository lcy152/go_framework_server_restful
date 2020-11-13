package impl

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	model "tumor_server/model"
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

func LoadAddInstitutionApplication(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	search := c.GetURLParam("search")
	status := c.GetURLParam("status")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)

	option := db.NewOptions()
	option.PageSize = pageSize
	option.PageIndex = pageIndex
	if search != "" {
		option.Match[db.OptUserName] = search
		option.Match[db.OptInstitutionName] = search
	}
	if status != "" {
		option.EQ[db.OptStatus] = status
	}
	sc := service.GetContainerInstance()
	insList, count, err := sc.DB.LoadAddInstitutionApplication(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseListHandler(c, count, insList)
}

func ApproveAddInstitutionApplication(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	op, err := sc.DB.GetAddInstitutionApplication(ctx, oid)
	CheckHandler(err, message.GetError)
	op.Institution.CreateTime = time.Now()
	err = sc.DB.AddInstitution(ctx, op.Institution)
	CheckHandler(err, message.AddError)
	uti := &model.UserToInstitution{
		ID:              primitive.NewObjectID(),
		Institution:     op.Institution.ID,
		InstitutionName: op.Institution.Name,
		Manager:         true,
		Type:            model.UserToInstitutionTypeDoctor,
		User:            op.User,
		UserName:        op.UserName,
	}
	err = sc.DB.AddUserToInstitution(ctx, uti)
	CheckHandler(err, message.AddError)
	op.Status = model.ApplicationStatusApprove
	op.OperateTime = time.Now()
	err = sc.DB.UpdateAddInstitutionApplication(ctx, op)
	CheckHandler(err, message.UpdateError)
	session.Commit()
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func RejectAddInstitutionApplication(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	op, err := sc.DB.GetAddInstitutionApplication(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	op.Status = model.ApplicationStatusReject
	op.OperateTime = time.Now()
	err = sc.DB.UpdateAddInstitutionApplication(context.TODO(), op)
	CheckHandler(err, message.UpdateError)
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}
