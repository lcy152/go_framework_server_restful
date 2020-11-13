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

func AddUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.User{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(data.Phone == "", message.RequestDataError)
	tn := time.Now()
	data.ID = NewUUID()
	data.CreateTime = tn
	data.LastModTime = tn
	sc := service.GetContainerInstance()
	err := sc.DB.AddUser(context.TODO(), data)
	CheckHandler(err, message.PhoneError)
	HttpReponseHandler(c, data)
}

func GetUser(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	HttpReponseHandler(c, userInfo.User)
}

func SearchUser(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)
	receiver := c.GetURLParam("receiver")
	option := db.NewOptions()
	option.EQ[db.OptHidden] = false
	option.Match[db.OptName] = receiver
	option.Match[db.OptPhone] = receiver
	option.Match[db.OptIDCard] = receiver
	option.PageIndex = pageIndex
	option.PageSize = pageSize
	sc := service.GetContainerInstance()
	lsit, count, err := sc.DB.LoadUser(context.TODO(), option)
	CheckHandler(err, message.GetListError)
	HttpReponseListHandler(c, count, lsit)
}

func EditUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ID{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	user, err := sc.DB.GetUser(context.TODO(), data.ID)
	CheckHandler(err, message.GetError)
	CheckHandler(!c.ParseBody(user), message.JsonParseError)
	CheckHandler(len(user.Photo) > 1024*1000, message.ImageSizeError)
	CheckHandler(len(user.Qrcode) > 1024*1000, message.ImageSizeError)
	err = service.UpdateUser(context.TODO(), user)
	CheckHandler(err, message.UpdateUserError)
	HttpReponseHandler(c, data)
}

func EditUserPassword(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(len(data.NewPassword) < 6, message.RequestDataError)
	sc := service.GetContainerInstance()
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	user, err := sc.DB.GetUser(context.TODO(), userInfo.User.ID)
	CheckHandler(err, message.FindUserError)
	CheckHandler(data.OldPassword != user.Password, message.PasswordError)
	err = sc.DB.UpdateUserPassword(context.TODO(), userInfo.User.ID, data.NewPassword)
	CheckHandler(err, message.UpdateUserError)
	HttpReponseHandler(c, nil)
}

func EditUserPhone(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		Phone string `json:"phone" bson:"phone"`
		Code  string `json:"code" bson:"code"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(len(data.Phone) < 6, message.RequestDataError)
	sc := service.GetContainerInstance()
	CheckHandler(!service.ShortMessageValidate(data.Phone, data.Code), message.ShortMessageValidateError)
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	err := sc.DB.UpdateUserPhone(context.TODO(), userInfo.User.ID, data.Phone)
	CheckHandler(err, message.UpdateUserError)
	service.DeleteUserTokenInfo(userInfo.User.ID.String())
	HttpReponseHandler(c, nil)
}

func UserFriendList(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.EQ[db.OptUser] = userInfo.User.ID
	fList, err := sc.DB.LoadUserToUser(context.TODO(), opt)
	CheckHandler(err, message.GetListError)
	HttpReponseHandler(c, fList)
}

func UserInstitutionList(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.EQ[db.OptUser] = userInfo.User.ID
	insList, err := sc.DB.LoadUserToInstitution(context.TODO(), opt)
	CheckHandler(err, message.GetListError)
	type E struct {
		model.Institution
		Type string `json:"type" bson:"type"`
		Job  string `json:"job" bson:"job"`
	}
	institutionList := []E{}
	for _, refIns := range insList {
		ins, err := sc.DB.GetInstitution(context.TODO(), refIns.Institution)
		CheckHandler(err, message.GetError)
		e := E{
			Institution: *ins,
			Type:        refIns.Type,
			Job:         refIns.Job,
		}
		institutionList = append(institutionList, e)
	}
	HttpReponseHandler(c, institutionList)
}

func AddFriend(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		User        primitive.ObjectID `json:"user" bson:"user"`
		Description string             `json:"description" bson:"description"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	user, err := sc.DB.GetUser(context.TODO(), data.User)
	CheckHandler(err, message.FindUserError)
	option := db.NewOptions()
	option.EQ[db.OptUser] = userInfo.User.ID
	option.EQ[db.OptFriend] = user.ID
	option.EQ[db.OptStatus] = model.ApplicationStatusWait
	uaList, _ := sc.DB.LoadUserApplication(context.TODO(), option)
	CheckHandler(len(uaList) != 0, message.RequestRepeatError)
	userApplication := &model.UserApplication{
		ID:          NewUUID(),
		Friend:      user.ID,
		FriendName:  user.Name,
		Status:      model.ApplicationStatusWait,
		User:        userInfo.User.ID,
		UserName:    userInfo.User.Name,
		Description: data.Description,
		CreateTime:  time.Now(),
	}
	err = sc.DB.AddUserApplication(context.TODO(), userApplication)
	CheckHandler(err, message.AddError)
	jsonStr, _ := json.Marshal(userApplication)
	service.SendFriendMessage(userApplication.ID.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func ApproveFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()

	userInfo := GetContextUserInfo(c)
	op := sc.DB.GetUserApplication(ctx, oid)
	CheckHandler(op == nil, message.GetError)
	CheckHandler(op.User != userInfo.User.ID, message.AuthorityError)
	op.Status = model.ApplicationStatusApprove
	op.OperateTime = time.Now()
	err = sc.DB.UpdateUserApplication(ctx, op)
	CheckHandler(err, message.UpdateError)

	utu1 := &model.UserToUser{
		ID:         NewUUID(),
		User:       op.User,
		UserName:   op.UserName,
		Friend:     op.Friend,
		FriendName: op.FriendName,
	}
	err = sc.DB.AddUserToUser(ctx, utu1)
	CheckHandler(err, message.AddError)
	utu2 := &model.UserToUser{
		ID:         NewUUID(),
		User:       op.Friend,
		UserName:   op.FriendName,
		Friend:     op.User,
		FriendName: op.UserName,
	}
	err = sc.DB.AddUserToUser(ctx, utu2)
	CheckHandler(err, message.AddError)

	session.Commit()
	jsonStr, _ := json.Marshal(op)
	service.SendFriendMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func RejectFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	op := sc.DB.GetUserApplication(context.TODO(), oid)
	CheckHandler(op == nil, message.GetError)
	CheckHandler(op.User != userInfo.User.ID, message.AuthorityError)
	op.Status = model.ApplicationStatusReject
	op.OperateTime = time.Now()
	err := sc.DB.UpdateUserApplication(context.TODO(), op)
	CheckHandler(err, message.UpdateError)
	jsonStr, _ := json.Marshal(op)
	service.SendFriendMessage(op.User.String(), string(jsonStr))
	HttpReponseHandler(c, nil)
}

func GetFriendApplicationList(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	opt := db.NewOptions()
	opt.EQ[db.OptUser] = userInfo.User.ID
	opList, err := sc.DB.LoadUserApplication(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, opList)
}

func DeleteFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()
	{
		opt := db.NewOptions()
		opt.EQ[db.OptUser] = userInfo.User.ID
		opt.EQ[db.OptFriend] = oid
		utuList, err := sc.DB.LoadUserToUser(ctx, opt)
		CheckHandler(err, message.GetListError)
		for _, u := range utuList {
			err := sc.DB.DeleteUserToUser(ctx, u.ID)
			CheckHandler(err, message.DeleteError)
		}
	}
	{
		opt := db.NewOptions()
		opt.EQ[db.OptUser] = oid
		opt.EQ[db.OptFriend] = userInfo.User.ID
		utuList, err := sc.DB.LoadUserToUser(ctx, opt)
		CheckHandler(err, message.GetListError)
		for _, u := range utuList {
			err := sc.DB.DeleteUserToUser(ctx, u.ID)
			CheckHandler(err, message.DeleteError)
		}
	}
	session.Commit()
	HttpReponseHandler(c, nil)
}

func GetUserDetailInstitution(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")
	userType := c.GetParam("user_type")
	oid, _ := primitive.ObjectIDFromHex(institutionID)
	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	utu, err := sc.DB.GetUserToInstitutionUserType(context.TODO(), oid, userInfo.User.ID, userType)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, utu)
}

func ChangeCurrentInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer session.Close()

	opt := db.NewOptions()
	opt.EQ[db.OptUser] = userInfo.User.ID
	urList, err := sc.DB.LoadUserToInstitution(context.TODO(), opt)
	for _, v := range urList {
		if v.ID == oid {
			v.Current = true
			err := sc.DB.UpdateUserToInstitution(ctx, v)
			CheckHandler(err, message.UpdateError)
		} else if v.Current {
			v.Current = false
			err := sc.DB.UpdateUserToInstitution(ctx, v)
			CheckHandler(err, message.UpdateError)
		}
	}
	session.Commit()
	HttpReponseHandler(c, nil)
}

func AuthDipperUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ValidateDipperUser{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	ur, _ := sc.DB.GetUserToInstitutionUserType(context.TODO(), data.Institution, userInfo.User.ID, "")
	CheckHandler(ur == nil, message.RequestDataError)
	service.ValidateDipperUserPublish(data)
	HttpReponseHandler(c, nil)
}
