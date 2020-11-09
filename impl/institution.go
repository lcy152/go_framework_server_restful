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
)

func AddInstitution(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Institution{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)

	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	tn := time.Now()
	if len(data.Manager) == 0 {
		data.Manager = []string{userInfo.User.Guid}
	}
	data.Guid = NewUUID()
	data.CreateTime = tn
	data.Creator = userInfo.User.Guid
	err = sc.DB.AddInstitution(ctx, data)
	CheckHandler(err, message.AddError)

	userRefIns := &model.UserRouter{
		Guid:            NewUUID(),
		InstitutionId:   data.Guid,
		InstitutionName: data.Name,
		UserGuid:        userInfo.User.Guid,
		Creator:         userInfo.User.Guid,
		LastOperator:    userInfo.User.Guid,
		CreatedTime:     time.Now(),
		LastModTime:     time.Now(),
	}
	opt := db.NewOptions()
	opt.Search[db.OptUserGuid] = userInfo.User.Guid
	opt.Search[db.OptCurrent] = true
	urList := sc.DB.LoadUserRouter(context.TODO(), opt)
	if len(urList) == 0 {
		userRefIns.IsCurrent = true
	}
	err = sc.DB.AddUserRouter(ctx, userRefIns)
	CheckHandler(err, message.AddError)

	ctx.CommitTransaction(ctx)
	sc.RabbitMQ.Connect()
	HttpReponseHandler(c, data)
}

func GetInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	ins := sc.DB.GetInstitution(context.TODO(), id)
	CheckHandler(ins == nil, message.GetError)
	HttpReponseHandler(c, ins)
}

func EditInstitution(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Institution{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	data = sc.DB.GetInstitution(context.TODO(), data.Guid)
	hasAuth := false
	for _, v := range data.Manager {
		if v == userInfo.User.Guid {
			hasAuth = true
			break
		}
	}
	if len(data.Manager) == 0 {
		data.Manager = append(data.Manager, userInfo.User.Guid)
	}
	CheckHandler(!hasAuth, message.AuthorityError)
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	err := sc.DB.UpdateInstitution(context.TODO(), data)
	CheckHandler(err, message.UpdateError)
	service.DeleteUserTokenInfo(data.Guid)
	HttpReponseHandler(c, data)
}

func DeleteInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	ins := sc.DB.GetInstitution(context.TODO(), id)
	hasAuth := false
	for _, v := range ins.Manager {
		if v == userInfo.User.Guid {
			hasAuth = true
			break
		}
	}
	CheckHandler(!hasAuth, message.AuthorityError)
	err := sc.DB.DeleteInstitution(context.TODO(), id)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func LoadInstitution(c *framework.Context) {
	defer PanicHandler(c)
	var data = struct {
		TimeStart  int64  `json:"time_start"`
		TimeEnd    int64  `json:"time_end"`
		PageIndex  int64  `json:"page_index"`
		PageSize   int64  `json:"page_size"`
		AscendSort bool   `json:"ascend_sort"`
		SortOption string `json:"sort_option"`
		Search     string `json:"search"`
	}{
		PageSize:   50,
		SortOption: "created_time",
	}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	option := db.NewOptions()
	option.TimeKey = db.OptCreateTime
	option.TimeStart = data.TimeStart
	option.TimeEnd = data.TimeEnd
	option.PageSize = data.PageSize
	option.PageIndex = data.PageIndex
	option.Ascend = data.AscendSort
	if data.Search != "" {
		option.Search[db.OptInstitutionId] = data.Search
		option.Search[db.OptInstitutionName] = data.Search
		option.Search[db.OptCode] = data.Search
		option.Search[db.OptAddress] = data.Search
		option.Regex[db.OptInstitutionId] = true
		option.Search[db.OptInstitutionName] = data.Search
		option.Search[db.OptCode] = data.Search
		option.Regex[db.OptAddress] = true
	}
	switch data.SortOption {
	case "created_time":
		option.Sort = append(option.Sort, db.OptCreateTime)
	case "guid":
		option.Sort = append(option.Sort, db.OptGuid)
	case "name":
		option.Sort = append(option.Sort, db.OptName)
	default:
		option.Sort = append(option.Sort, db.OptCreateTime)
	}
	sc := service.GetContainerInstance()
	insList, _, err := sc.DB.LoadInstitution(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, insList)
}

func LoadUserInstitution(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	insList := []*model.Institution{}
	opt := db.NewOptions()
	opt.Search[db.OptUserGuid] = userInfo.User.Guid
	urList := sc.DB.LoadUserRouter(context.TODO(), opt)
	for _, v := range urList {
		refIns := sc.DB.GetUserRouter(context.TODO(), v.InstitutionId)
		if refIns == nil {
			continue
		}
		ins := sc.DB.GetInstitution(context.TODO(), refIns.Guid)
		if ins != nil {
			insList = append(insList, ins)
		}
	}
	HttpReponseHandler(c, insList)
}

func ApplyInstitution(c *framework.Context) {
	institutionID := c.GetParam("institution_id")
	description := c.GetParam("description")
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	institution := sc.DB.GetInstitution(context.TODO(), institutionID)
	CheckHandler(institution == nil, message.GetError)
	option := db.NewOptions()
	option.Search[db.OptCreator] = userInfo.User.Guid
	option.Search[db.OptInstitutionId] = institution.Guid
	uaList, _, _ := sc.DB.LoadUserApplication(context.TODO(), option)
	for _, v := range uaList {
		CheckHandler(v.Status == model.ApplicationStatusWait, message.RequestRepeatError)
	}
	userApplication := &model.UserApplication{
		Guid:          NewUUID(),
		InstitutionId: institution.Guid,
		Type:          model.ApplicationTypeApplyInstitution,
		Status:        model.ApplicationStatusWait,
		Creator:       userInfo.User.Guid,
		CreateTime:    time.Now(),
		Description:   description,
	}
	err := sc.DB.AddUserApplication(context.TODO(), userApplication)
	CheckHandler(err, message.AddError)
	for _, v := range institution.Manager {
		jsonStr, _ := json.Marshal(userApplication)
		service.SendInstitutionMessage(v, string(jsonStr))
	}
	HttpReponseHandler(c, nil)
}

func InstitutionUserList(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")

	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.Search[db.OptInstitutionId] = institutionID
	urList := sc.DB.LoadUserRouter(context.TODO(), opt)
	CheckHandler(len(urList) == 0, message.AuthorityError)
	list := []*model.User{}
	in := true
	for _, v := range urList {
		user := sc.DB.GetUser(context.TODO(), v.UserGuid)
		if user == nil {
			continue
		}
		if user.Guid == userInfo.User.Guid {
			in = false
		}
		list = append(list, user)
	}
	CheckHandler(in, message.GetError)
	HttpReponseHandler(c, list)
}

func ApproveInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")

	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	op := sc.DB.GetUserApplication(ctx, id)
	CheckHandler(op == nil, message.GetError)

	institution := sc.DB.GetInstitution(context.TODO(), op.InstitutionId)
	CheckHandler(institution == nil, message.GetError)
	hasAuth := false
	for _, v := range institution.Manager {
		if userInfo.User.Guid == v {
			hasAuth = true
			break
		}
	}
	CheckHandler(!hasAuth, message.AuthorityError)

	user := sc.DB.GetUser(context.TODO(), op.Creator)
	CheckHandler(user == nil, message.FindUserError)

	userRefIns := &model.UserRouter{
		Guid:            NewUUID(),
		InstitutionId:   institution.Guid,
		InstitutionName: institution.Name,
		Creator:         user.Guid,
		LastOperator:    user.Guid,
		CreatedTime:     time.Now(),
		LastModTime:     time.Now(),
	}
	opt := db.NewOptions()
	opt.Search[db.OptUserGuid] = user.Guid
	opt.Search[db.OptCurrent] = true
	urList := sc.DB.LoadUserRouter(context.TODO(), opt)
	if len(urList) == 0 {
		userRefIns.IsCurrent = true
	}
	err = sc.DB.AddUserRouter(ctx, userRefIns)
	CheckHandler(err, message.AddError)

	op.Status = model.ApplicationStatusApprove
	err = sc.DB.UpdateUserApplication(ctx, op)
	CheckHandler(err, message.UpdateError)

	ctx.CommitTransaction(ctx)
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.Creator, string(jsonStr))
	HttpReponseHandler(c, nil)
}

func RejectInstitution(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")

	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	op := sc.DB.GetUserApplication(context.TODO(), id)
	CheckHandler(op == nil, message.GetError)
	institution := sc.DB.GetInstitution(context.TODO(), op.InstitutionId)
	CheckHandler(institution == nil, message.GetError)
	hasAuth := false
	for _, v := range institution.Manager {
		if userInfo.User.Guid == v {
			hasAuth = true
			break
		}
	}
	CheckHandler(!hasAuth, message.AuthorityError)
	op.Status = model.ApplicationStatusReject
	err := sc.DB.UpdateUserApplication(context.TODO(), op)
	CheckHandler(err, message.UpdateError)
	jsonStr, _ := json.Marshal(op)
	service.SendInstitutionMessage(op.Creator, string(jsonStr))
	HttpReponseHandler(c, nil)
}

func GetInstitutionApply(c *framework.Context) {
	defer PanicHandler(c)
	institutionID := c.GetParam("institution_id")
	state := c.GetParam("state")
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	ins := sc.DB.GetInstitution(context.TODO(), institutionID)
	CheckHandler(ins == nil, message.GetError)
	for _, v := range ins.Manager {
		if v == userInfo.User.Guid {
			opt := db.NewOptions()
			opt.Search[db.OptInstitutionId] = institutionID
			opt.Search[db.OptType] = model.ApplicationTypeApplyInstitution
			opt.Search[db.OptStatus] = state
			opList, _, err := sc.DB.LoadUserApplication(context.TODO(), opt)
			CheckHandler(err, message.GetError)
			HttpReponseHandler(c, opList)
			return
		}
	}
	CheckHandler(true, message.AuthorityError)
}
