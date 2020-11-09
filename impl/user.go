package impl

import (
	"context"
	"encoding/json"
	"time"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	uuid "github.com/satori/go.uuid"
)

func AddUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.User{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(data.Phone == "", message.RequestDataError)
	tn := time.Now()
	data.Guid = uuid.Must(uuid.NewV4(), nil).String()
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
	var data = &struct {
		Search     string `json:"search"`
		PageIndex  int64  `json:"page_index"`
		PageSize   int64  `json:"page_size"`
		AscendSort bool   `json:"ascend_sort"`
		SortOption string `json:"sort_option"`
	}{PageIndex: 0, PageSize: 20}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	option := db.NewOptions()
	option.TimeKey = db.OptCreateTime
	option.Search[db.OptDisable] = false
	option.Search[db.OptName] = data.Search
	option.Search[db.OptPhone] = data.Search
	option.Search[db.OptIDCard] = data.Search
	option.Regex[db.OptName] = true
	option.Regex[db.OptPhone] = true
	option.Regex[db.OptIDCard] = true
	option.PageIndex = data.PageIndex
	option.PageSize = data.PageSize
	option.Ascend = data.AscendSort
	switch data.SortOption {
	case "created_time":
		option.Sort = append(option.Sort, db.OptCreateTime)
	case "guid":
		option.Sort = append(option.Sort, db.OptGuid)
	case "name":
		option.Sort = append(option.Sort, db.OptName)
	case "phone":
		option.Sort = append(option.Sort, db.OptPhone)
	default:
		option.Sort = append(option.Sort, db.OptCreateTime)
	}
	sc := service.GetContainerInstance()
	lsit, _, err := sc.DB.LoadUser(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, lsit)
}

func EditUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Guid{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	user := sc.DB.GetUser(context.TODO(), data.Guid)
	CheckHandler(!c.ParseBody(user), message.JsonParseError)
	CheckHandler(len(user.Photo) > 1024*1000, message.ImageSizeError)
	CheckHandler(len(user.Qrcode) > 1024*1000, message.ImageSizeError)
	err := service.UpdateUser(context.TODO(), user)
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
	user := sc.DB.GetUser(context.TODO(), userInfo.User.Guid)
	CheckHandler(user == nil, message.FindUserError)
	CheckHandler(data.OldPassword != user.Password, message.PasswordError)
	err := sc.DB.UpdateUserPassword(context.TODO(), userInfo.User.Guid, data.NewPassword)
	CheckHandler(err, message.UpdateUserError)
	HttpReponseHandler(c, nil)
}

func EditUserPhone(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		Phone string `json:"phone" bson:"phone"`
		Code  int64  `json:"code" bson:"code"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(len(data.Phone) < 6, message.RequestDataError)
	sc := service.GetContainerInstance()
	CheckHandler(!service.ShortMessageValidate(data.Phone, data.Code), message.ShortMessageValidateError)
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	err := sc.DB.UpdateUserPhone(context.TODO(), userInfo.User.Guid, data.Phone)
	CheckHandler(err, message.UpdateUserError)
	service.DeleteUserTokenInfo(userInfo.User.Guid)
	HttpReponseHandler(c, nil)
}

func GetUserFriendList(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	fList := []*model.User{}
	sc := service.GetContainerInstance()
	for _, v := range userInfo.User.FriendList {
		user := sc.DB.GetUser(context.TODO(), v)
		if user != nil {
			fList = append(fList, user)
		}
	}
	HttpReponseHandler(c, fList)
}

func AddFriend(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		UserGuid    string `json:"user_guid" bson:"user_guid"`
		Description string `json:"description" bson:"description"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	user := sc.DB.GetUser(context.TODO(), data.UserGuid)
	CheckHandler(user == nil, message.FindUserError)
	option := db.NewOptions()
	option.Search[db.OptCreator] = userInfo.User.Guid
	option.Search[db.OptUserGuid] = user.Guid
	option.Search[db.OptStatus] = model.ApplicationStatusWait
	uaList, _, _ := sc.DB.LoadUserApplication(context.TODO(), option)
	CheckHandler(len(uaList) != 0, message.RequestRepeatError)
	userApplication := &model.UserApplication{
		Guid:        NewUUID(),
		UserGuid:    user.Guid,
		Type:        model.ApplicationTypeApplyFriend,
		Status:      model.ApplicationStatusWait,
		Creator:     userInfo.User.Guid,
		CreateTime:  time.Now(),
		Description: data.Description,
	}
	err := sc.DB.AddUserApplication(context.TODO(), userApplication)
	CheckHandler(err, message.AddError)
	jsonStr, _ := json.Marshal(userApplication)
	service.SendFriendMessage(userApplication.UserGuid, string(jsonStr))
	HttpReponseHandler(c, nil)
}

func ApproveFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	userInfo := GetContextUserInfo(c)
	op := sc.DB.GetUserApplication(ctx, id)
	CheckHandler(op == nil, message.GetError)
	CheckHandler(op.UserGuid != userInfo.User.Guid, message.AuthorityError)
	user := userInfo.User
	user.FriendList = append(user.FriendList, op.Creator)
	err = service.UpdateUser(ctx, user)
	CheckHandler(err, message.UpdateError)
	op.Status = model.ApplicationStatusApprove
	err = sc.DB.UpdateUserApplication(ctx, op)
	CheckHandler(err, message.UpdateError)
	ctx.CommitTransaction(ctx)
	jsonStr, _ := json.Marshal(op)
	service.SendFriendMessage(op.Creator, string(jsonStr))
	HttpReponseHandler(c, nil)
}

func RejectFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	op := sc.DB.GetUserApplication(context.TODO(), id)
	CheckHandler(op == nil, message.GetError)
	CheckHandler(op.UserGuid != userInfo.User.Guid, message.AuthorityError)
	op.Status = model.ApplicationStatusReject
	err := sc.DB.UpdateUserApplication(context.TODO(), op)
	CheckHandler(err, message.UpdateError)
	HttpReponseHandler(c, nil)
}

func GetFriendApplicationList(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	opt := db.NewOptions()
	opt.Search[db.OptUserGuid] = userInfo.User.Guid
	opt.Search[db.OptType] = model.ApplicationTypeApplyFriend
	opList, _, err := sc.DB.LoadUserApplication(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, opList)
}

func DeleteFriend(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	friendList := []string{}
	for _, v := range userInfo.User.FriendList {
		if v == id {
			continue
		}
		friendList = append(friendList, v)
	}
	err = service.UpdateUser(ctx, userInfo.User)
	CheckHandler(err, message.UpdateError)

	friend := sc.DB.GetUser(ctx, id)
	CheckHandler(friend == nil, message.GetError)

	friendList2 := []string{}
	for _, v := range friend.FriendList {
		if v == userInfo.User.Guid {
			continue
		}
		friendList2 = append(friendList2, v)
	}
	err = service.UpdateUser(ctx, friend)
	CheckHandler(err, message.UpdateError)
	ctx.CommitTransaction(ctx)
	HttpReponseHandler(c, nil)
}

func ChangeUserInstitution(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		Guid string `json:"guid" bson:"guid"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	noIns := true
	opt := db.NewOptions()
	opt.Search[db.OptUserGuid] = userInfo.User.Guid
	urList := sc.DB.LoadUserRouter(context.TODO(), opt)
	for _, v := range urList {
		refIns := sc.DB.GetUserRouter(ctx, v.InstitutionId)
		if refIns == nil {
			continue
		}
		if refIns.Guid == data.Guid {
			refIns.IsCurrent = true
			noIns = false
		} else {
			refIns.IsCurrent = false
		}
		err := sc.DB.UpdateUserRouter(ctx, refIns)
		CheckHandler(err, message.UpdateError)
	}
	CheckHandler(noIns, message.GetError)
	ctx.CommitTransaction(ctx)
	HttpReponseHandler(c, nil)
}

func AuthDipperUser(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ValidateDipperUser{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)

	sc := service.GetContainerInstance()
	userInfo := GetContextUserInfo(c)
	ur := sc.DB.GetUserRouterByTumorUser(context.TODO(), data.InstitutionId, userInfo.User.Guid)
	CheckHandler(ur == nil, message.RequestDataError)
	service.ValidateDipperUserPublish(data)
	HttpReponseHandler(c, nil)
}
