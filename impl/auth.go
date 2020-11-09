package impl

import (
	"context"
	"log"
	"strconv"
	"time"
	framework "tumor_server/framework"
	message "tumor_server/message"
	model "tumor_server/model"
	service "tumor_server/service"

	uuid "github.com/satori/go.uuid"
)

func Login(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	data := struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}{}
	CheckHandler(!c.ParseBody(&data), message.JsonParseError)
	CheckHandler(data.Phone == "" || data.Password == "", message.RequestDataError)

	user := sc.DB.GetUserByPhone(context.TODO(), data.Phone)
	CheckHandler(user == nil, message.FindUserError)
	CheckHandler(data.Password != user.Password, message.PasswordError)
	tn := time.Now().Unix()
	ti := &service.Token{
		Guid:      user.Guid,
		LoginTime: tn,
	}
	token, err := service.MarshalToken(ti)
	CheckHandler(err, message.TokenGenerateError)
	userInfo := &service.UserTokenInfo{
		User:      user,
		Token:     token,
		LoginTime: tn,
		IP:        c.Req.Host,
	}
	err = service.AddUserTokenInfo(userInfo)
	if err != nil {
		log.Println("error: AddRedisSession", err)
	}
	err = service.GetContainerInstance().DB.UpdateUserToken(context.TODO(), user.Guid, token)
	CheckHandler(err, message.UpdateUserError)
	service.AddUserRecord(user.Guid, model.UserOperationLogin, "", c.Req.Host)
	c.SetHeader("Authorization", token)
	HttpReponseHandler(c, user)
}

func LoginCode(c *framework.Context) {
	defer PanicHandler(c)
	data := struct {
		Phone string `json:"phone"`
		Code  int64  `json:"code"`
	}{}
	CheckHandler(!c.ParseBody(&data), message.JsonParseError)
	CheckHandler(data.Phone == "" || data.Code == 0, message.RequestDataError)
	CheckHandler(!service.ShortMessageValidate(data.Phone, data.Code), message.ShortMessageValidateError)

	sc := service.GetContainerInstance()
	session, ctx, err := sc.DB.StartSession()
	CheckHandler(err, message.HttpError)
	defer func() {
		sc.DB.EndSession(session)
		ctx.AbortTransaction(ctx)
	}()
	ctx.StartTransaction()

	user := sc.DB.GetUserByPhone(ctx, data.Phone)
	if user == nil {
		user = &model.User{}
		user.Guid = NewUUID()
		user.Phone = data.Phone
		user.Password = NewUUID()
		err := sc.DB.AddUser(ctx, user)
		CheckHandler(err, message.RegisterError)
	}
	tn := time.Now().Unix()
	ti := &service.Token{
		Guid:      user.Guid,
		LoginTime: tn,
	}
	token, err := service.MarshalToken(ti)
	CheckHandler(err, message.TokenGenerateError)
	userInfo := &service.UserTokenInfo{
		User:      user,
		Token:     token,
		LoginTime: tn,
		IP:        c.Req.Host,
	}
	err = service.AddUserTokenInfo(userInfo)
	if err != nil {
		log.Println("error: AddRedisSession", err)
	}
	service.AddUserRecord(user.Guid, model.UserOperationLogin, "", c.Req.Host)
	err = service.StoreShortMessage(data.Phone, 0)
	if err != nil {
		log.Println("error: StoreShortMessage", err)
	}
	user.Token = token
	err = sc.DB.UpdateUserToken(ctx, user.Guid, token)
	CheckHandler(err, message.UpdateUserError)
	ctx.CommitTransaction(ctx)
	c.SetHeader("Authorization", token)
	HttpReponseHandler(c, user)
}

func Register(c *framework.Context) {
	defer PanicHandler(c)
	data := struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Code     int64  `json:"code"`
	}{}
	CheckHandler(!c.ParseBody(&data), message.JsonParseError)
	CheckHandler(data.Phone == "" || data.Password == "" || data.Code == 0, message.RequestDataError)
	CheckHandler(!service.ShortMessageValidate(data.Phone, data.Code), message.ShortMessageValidateError)

	sc := service.GetContainerInstance()
	user := sc.DB.GetUserByPhone(context.TODO(), data.Phone)
	CheckHandler(user != nil, message.UserExistError)
	user = &model.User{}
	user.Guid = uuid.Must(uuid.NewV4(), nil).String()
	user.Phone = data.Phone
	user.Password = uuid.Must(uuid.NewV4(), nil).String()
	err := sc.DB.AddUser(context.TODO(), user)
	CheckHandler(err, message.RegisterError)
	err = service.StoreShortMessage(data.Phone, 0)
	if err != nil {
		log.Println("error: StoreShortMessage", err)
	}
	HttpReponseHandler(c, user)
}

func GetShortMessageCode(c *framework.Context) {
	defer PanicHandler(c)
	phone := c.GetParam("phone")
	CheckHandler(phone == "", message.RequestDataError)
	code := service.GenerateCode(4)
	text := "您的验证码是：" + strconv.Itoa(int(code)) + "。请不要把验证码泄露给其他人。"
	// i := SendMessageToIhuyi(phone, text)
	i := &IhuyiResponse{
		Code: 2,
		Msg:  text,
	}
	CheckHandler(i == nil, message.HttpError)
	CheckHandler(i.Code != 2, message.HttpError)
	err := service.StoreShortMessage(phone, code)
	if err != nil {
		CheckHandler(true, err.Error())
	}
	HttpReponseHandler(c, i)
}

func LogOut(c *framework.Context) {
	defer PanicHandler(c)
	userInfo, _ := service.UnmarshalToken(c.GetAuthorization())
	service.DeleteUserTokenInfo(userInfo.Guid)
	service.AddUserRecord(userInfo.Guid, model.UserOperationLogout, "", c.Req.Host)
	_ = service.GetContainerInstance().DB.UpdateUserToken(context.TODO(), userInfo.Guid, "")
	HttpReponseHandler(c, nil)
}
