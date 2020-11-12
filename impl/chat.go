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

func AddSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ChatMessage{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	data.Sender = userInfo.User.ID
	data.SenderName = userInfo.User.Name
	data.ID = NewUUID()
	data.CreateTime = time.Now()
	err := sc.DB.AddSingleChat(context.TODO(), data)
	CheckHandler(err, message.AddError)
	jsonStr, _ := json.Marshal(data)
	service.SendSingleMessage(data.Receiver.String(), string(jsonStr))
	HttpReponseHandler(c, data)
}

func DeleteSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	info := sc.DB.GetSingleChat(context.TODO(), oid)
	CheckHandler(info == nil, message.GetError)
	CheckHandler(time.Now().Unix()-info.CreateTime.Unix() > 120, message.ExpiredError)
	err := sc.DB.DeleteSingleChat(context.TODO(), oid)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func LoadSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)
	receiver := c.GetURLParam("receiver")

	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.PageIndex = pageIndex
	opt.PageSize = pageSize

	query1 := map[db.OptionKey]interface{}{}
	query1[db.OptSender] = userInfo.User.ID
	query1[db.OptReceiver] = receiver
	query2 := map[db.OptionKey]interface{}{}
	query2[db.OptSender] = receiver
	query2[db.OptReceiver] = userInfo.User.ID
	opt.OR = append(opt.OR, query1, query2)

	sgList, err := sc.DB.LoadSingleChat(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, sgList)
}

func AddGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ChatMessage{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	data.CreateTime = time.Now()
	data.Sender = userInfo.User.ID
	data.SenderName = userInfo.User.Name
	data.ID = NewUUID()
	err := sc.DB.AddGroupChat(context.TODO(), data)
	CheckHandler(err, message.AddError)
	utuList, err := sc.DB.LoadUserToUserGroupUser(context.TODO(), data.Group)
	jsonStr, _ := json.Marshal(data)
	for _, v := range utuList {
		if data.Sender == v.User {
			continue
		}
		service.SendGroupMessage(v.User.String(), string(jsonStr))
	}
	HttpReponseHandler(c, data)
}

func DeleteGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	info, err := sc.DB.GetGroupChat(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	CheckHandler(time.Now().Unix()-info.CreateTime.Unix() > 120, message.ExpiredError)
	err = sc.DB.DeleteGroupChat(context.TODO(), oid)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func LoadGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)
	search := c.GetURLParam("search")
	group := c.GetURLParam("group")

	opt := db.NewOptions()
	opt.EQ[db.OptGroup] = group
	if search != "" {
		opt.EQ[db.OptType] = "text"
		opt.Match[db.OptData] = search
	}
	opt.PageIndex = pageIndex
	opt.PageSize = pageSize
	sc := service.GetContainerInstance()
	sgList, err := sc.DB.LoadGroupChat(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, sgList)
}

func GetUserChatHistory(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	ugList, err := sc.DB.LoadUserToUserGroupUser(context.TODO(), userInfo.User.ID)
	CheckHandler(err, message.GetError)
	friendList, err := sc.DB.LoadUserToUserUser(context.TODO(), userInfo.User.ID)
	var Message []*model.ChatMessage
	for _, v := range friendList {
		opt := db.NewOptions()
		opt.PageIndex = 0
		opt.PageSize = 1
		query1 := map[db.OptionKey]interface{}{}
		query1[db.OptSender] = userInfo.User.ID
		query1[db.OptReceiver] = v
		query2 := map[db.OptionKey]interface{}{}
		query2[db.OptSender] = v
		query2[db.OptReceiver] = userInfo.User.ID
		opt.OR = append(opt.OR, query1, query2)
		opt.Sort = []db.SortOption{{Key: db.OptID, Ascend: false}}
		sgList, _ := sc.DB.LoadSingleChat(context.TODO(), opt)
		if len(sgList) > 0 {
			Message = append(Message, sgList[0])
		}
	}
	for _, v := range ugList {
		opt := db.NewOptions()
		opt.EQ[db.OptGroup] = v.ID
		opt.PageIndex = 0
		opt.PageSize = 1
		opt.Sort = []db.SortOption{{Key: db.OptID, Ascend: false}}
		megList, _ := sc.DB.LoadGroupChat(context.TODO(), opt)
		if len(megList) > 0 {
			Message = append(Message, megList[0])
		}
	}
	HttpReponseHandler(c, Message)
}
