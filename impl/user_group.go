package impl

import (
	"context"
	"time"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	uuid "github.com/satori/go.uuid"
)

func AddUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		Manager []string `json:"manager" bson:"manager"`
		Member  []string `json:"member" bson:"member"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(len(data.Member) == 0, message.RequestDataError)
	tn := time.Now()
	ug := &model.UserGroup{}
	ug.Guid = uuid.Must(uuid.NewV4(), nil).String()
	ug.CreateTime = tn
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	data.Manager = append(data.Manager, userInfo.User.Guid)
	memberMap := make(map[string]*model.User)
	memberMap[userInfo.User.Guid] = userInfo.User
	managerMap := make(map[string]*model.User)
	managerMap[userInfo.User.Guid] = userInfo.User
	for _, v := range data.Manager {
		if _, ok := managerMap[v]; !ok {
			user := sc.DB.GetUser(context.TODO(), v)
			if user != nil {
				managerMap[user.Guid] = user
				memberMap[user.Guid] = user
			}
		}
	}
	for _, v := range data.Member {
		if _, ok := memberMap[v]; !ok {
			user := sc.DB.GetUser(context.TODO(), v)
			if user != nil {
				memberMap[user.Guid] = user
			}
		}
	}
	for _, v := range managerMap {
		ug.Manager = append(ug.Manager, v.Guid)
	}
	for _, v := range memberMap {
		ug.Member = append(ug.Member, v.Guid)
	}
	err := sc.DB.AddUserGroup(context.TODO(), ug)
	CheckHandler(err, message.PhoneError)
	HttpReponseHandler(c, ug)
}

func GetUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	group := sc.DB.GetUserGroup(context.TODO(), id)
	CheckHandler(group == nil, message.GetError)
	for i, v := range group.Member {
		user := sc.DB.GetUser(context.TODO(), v)
		if user != nil {
			group.Member[i] = user.Guid
		}
	}
	HttpReponseHandler(c, group)
}

func DeleteUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	err := sc.DB.DeleteUserGroup(context.TODO(), id)
	CheckHandler(err, message.DeleteError)
	HttpReponseHandler(c, nil)
}

func EditUserGroup(c *framework.Context) {
	defer PanicHandler(c)
	data := &struct {
		Guid    string   `json:"guid"`
		Manager []string `json:"manager" bson:"manager"`
		Member  []string `json:"member" bson:"member"`
	}{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	CheckHandler(len(data.Member) == 0, message.RequestDataError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	ug := sc.DB.GetUserGroup(context.TODO(), data.Guid)
	CheckHandler(ug == nil, message.GetError)
	managerFlag := false
	for _, v := range ug.Manager {
		if userInfo.User.Guid == v {
			managerFlag = true
			break
		}
	}
	CheckHandler(!managerFlag, message.AuthorityError)
	data.Manager = append(data.Manager, userInfo.User.Guid)
	memberMap := make(map[string]*model.User)
	memberMap[userInfo.User.Guid] = userInfo.User
	managerMap := make(map[string]*model.User)
	managerMap[userInfo.User.Guid] = userInfo.User
	for _, v := range data.Manager {
		if _, ok := managerMap[v]; !ok {
			user := sc.DB.GetUser(context.TODO(), v)
			if user != nil {
				managerMap[user.Guid] = user
				memberMap[user.Guid] = user
			}
		}
	}
	for _, v := range data.Member {
		if _, ok := memberMap[v]; !ok {
			user := sc.DB.GetUser(context.TODO(), v)
			if user != nil {
				memberMap[user.Guid] = user
			}
		}
	}
	ug.Manager = []string{}
	ug.Member = []string{}
	for _, v := range managerMap {
		ug.Manager = append(ug.Manager, v.Guid)
	}
	for _, v := range memberMap {
		ug.Member = append(ug.Member, v.Guid)
	}
	err := sc.DB.UpdateUserGroup(context.TODO(), ug)
	CheckHandler(err, message.PhoneError)
	HttpReponseHandler(c, data)
}

func GetUserGroupList(c *framework.Context) {
	defer PanicHandler(c)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.Search[db.OptMember] = userInfo.User.Guid
	ugList, _, err := sc.DB.LoadUserGroup(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, ugList)
}
