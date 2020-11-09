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

func AddSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ChatMessage{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(data.ReceiverGuid == "" || data.Data == "", message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	data.CreateTime = time.Now()
	data.SenderGuid = userInfo.User.Guid
	data.SenderName = userInfo.User.Name
	data.Guid = uuid.Must(uuid.NewV4(), nil).String()
	err := sc.DB.AddSingleChat(context.TODO(), data)
	CheckHandler(err, message.AddError)
	jsonStr, _ := json.Marshal(data)
	service.SendSingleMessage(data.ReceiverGuid, string(jsonStr))
	HttpReponseHandler(c, data)
}

func DeleteSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	info := sc.DB.GetSingleChat(context.TODO(), id)
	CheckHandler(info == nil, message.GetError)
	CheckHandler(time.Now().Unix()-info.CreateTime.Unix() > 120, message.ExpiredError)
	err := sc.DB.DeleteSingleChat(context.TODO(), id)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func LoadSingleChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	var data = &struct {
		Receiver  string `json:"receiver"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
		Ascend    bool   `json:"ascend_sort"`
		Search    string `json:"search"`
	}{
		PageSize: 20,
	}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.PageIndex = data.PageIndex
	opt.PageSize = data.PageSize
	opt.Ascend = data.Ascend
	sgList, err := sc.DB.LoadSingleChat(context.TODO(), opt, data.Receiver, userInfo.User.Guid, data.Search)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, sgList)
}

func AddGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.ChatMessage{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	group := sc.DB.GetUserGroup(context.TODO(), data.GroupGuid)
	CheckHandler(group == nil, message.GetError)
	hasAuth := false
	for _, v := range group.Member {
		if v == userInfo.User.Guid {
			hasAuth = true
			break
		}
	}
	CheckHandler(!hasAuth, message.AuthorityError)
	data.CreateTime = time.Now()
	data.SenderGuid = userInfo.User.Guid
	data.SenderName = userInfo.User.Name
	data.Guid = uuid.Must(uuid.NewV4(), nil).String()
	err := sc.DB.AddGroupChat(context.TODO(), data)
	CheckHandler(err, message.AddError)
	jsonStr, _ := json.Marshal(data)
	for _, v := range group.Member {
		if data.SenderGuid == v {
			continue
		}
		service.SendGroupMessage(v, string(jsonStr))
	}
	HttpReponseHandler(c, data)
}

func DeleteGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	info := sc.DB.GetGroupChat(context.TODO(), id)
	CheckHandler(info == nil, message.GetError)
	CheckHandler(time.Now().Unix()-info.CreateTime.Unix() > 120, message.ExpiredError)
	err := sc.DB.DeleteGroupChat(context.TODO(), id)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func LoadGroupChatMessage(c *framework.Context) {
	defer PanicHandler(c)
	var data = &struct {
		GroupGuid string `json:"group_guid"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
		Ascend    bool   `json:"ascend_sort"`
		Search    string `json:"search"`
	}{
		PageSize: 20,
	}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	opt := db.NewOptions()
	opt.Search[db.OptGroupId] = data.GroupGuid
	if data.Search != "" {
		opt.Search[db.OptData] = data.Search
	}
	opt.PageIndex = data.PageIndex
	opt.PageSize = data.PageSize
	opt.Ascend = data.Ascend
	opt.Sort = append(opt.Sort, db.OptCreateTime)
	sc := service.GetContainerInstance()
	sgList, err := sc.DB.LoadGroupChat(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, sgList)
}

func GetUserChatHistory(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.Search[db.OptMember] = userInfo.User.Guid
	ugList, _, err := sc.DB.LoadUserGroup(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	friendList := []*model.User{}
	for _, v := range userInfo.User.FriendList {
		user := sc.DB.GetUser(context.TODO(), v)
		if user != nil {
			friendList = append(friendList, user)
		}
	}
	var Message []*model.ChatMessage
	for _, v := range friendList {
		opt := db.NewOptions()
		opt.PageIndex = 0
		opt.PageSize = 1
		opt.Ascend = false
		megList, err := sc.DB.LoadSingleChat(context.TODO(), opt, v.Guid, userInfo.User.Guid, "")
		if err != nil || len(megList) == 0 {
			continue
		}
		Message = append(Message, megList[0])
	}
	for _, v := range ugList {
		opt := db.NewOptions()
		opt.Search[db.OptGroupId] = v.Guid
		opt.PageIndex = 0
		opt.PageSize = 1
		opt.Ascend = false
		opt.Sort = append(opt.Sort, db.OptCreateTime)
		megList, err := sc.DB.LoadGroupChat(context.TODO(), opt)
		if err != nil || len(megList) == 0 {
			continue
		}
		Message = append(Message, megList[0])
	}
	HttpReponseHandler(c, Message)
}
