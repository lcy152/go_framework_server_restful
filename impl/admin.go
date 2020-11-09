package impl

import (
	"context"
	"regexp"
	"time"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"
)

func GetAppConfig(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	lsit, err := sc.DB.LoadAppConfig(context.TODO())
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, lsit)
}

func GetUserList(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	var data = struct {
		TimeStart  int64  `json:"time_start"`
		TimeEnd    int64  `json:"time_end"`
		PageIndex  int64  `json:"page_index"`
		PageSize   int64  `json:"page_size"`
		AscendSort bool   `json:"ascend_sort"`
		SortOption string `json:"sort_option"`
		Search     string `json:"search"`
		Type       string `json:"type"`
	}{
		PageSize:   50,
		SortOption: "last_mod_time",
	}
	option := db.NewOptions()
	option.TimeKey = db.OptCreateTime
	option.TimeStart = data.TimeStart
	option.TimeEnd = data.TimeEnd
	option.PageSize = data.PageSize
	option.PageIndex = data.PageIndex
	option.Ascend = data.AscendSort
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	if data.Type != "" {
		option.Search[db.OptType] = data.Type
	}
	if data.Search != "" {
		option.Search[db.OptName] = data.Search
		option.Search[db.OptPhone] = data.Search
		option.Search[db.OptIDCard] = data.Search
		option.Regex[db.OptName] = true
		option.Regex[db.OptPhone] = true
		option.Regex[db.OptIDCard] = true
	}
	switch data.SortOption {
	case "created_time":
		option.Sort = append(option.Sort, db.OptCreateTime)
	case "last_mod_time":
		option.Sort = append(option.Sort, db.OptLastModTime)
	case "guid":
		option.Sort = append(option.Sort, db.OptGuid)
	case "name":
		option.Sort = append(option.Sort, db.OptName)
	default:
		option.Sort = append(option.Sort, db.OptCreateTime)
	}
	lsit, _, err := sc.DB.LoadUser(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, lsit)
}

func GetUserActivity(c *framework.Context) {
	defer PanicHandler(c)
	request := struct {
		UserID    string `json:"user_id"`
		TimeStart int64  `json:"time_start"`
		TimeEnd   int64  `json:"time_end"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
		Ascend    bool   `json:"ascend"`
	}{Ascend: true, PageIndex: 0, PageSize: 100000}
	CheckHandler(!c.ParseBody(&request), message.JsonParseError)
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	option := db.NewOptions()
	userTarget := sc.DB.GetUser(context.TODO(), request.UserID)
	CheckHandler(userTarget == nil, message.FindUserError)
	CheckHandler(userInfo.User.Guid != model.AdminGuid, message.FindUserError)
	option.Search[db.OptUserGuid] = userTarget.Guid
	option.TimeStart = request.TimeStart
	option.TimeKey = db.OptCreateTime
	option.TimeEnd = request.TimeEnd
	option.PageSize = request.PageSize
	option.PageIndex = request.PageIndex
	option.Ascend = request.Ascend
	option.Sort = append(option.Sort, db.OptCreateTime)
	list, count, err := sc.DB.LoadUserOperation(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseExtraListHandler(c, count, userTarget.Name, list)
}

func GetUserLoginRecord(c *framework.Context) {
	defer PanicHandler(c)
	const LogOut = model.UserOperationLogout
	const LogIn = model.UserOperationLogin
	request := struct {
		TimeStart int64  `json:"start_time"`
		TimeEnd   int64  `json:"end_time"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
		SearchKey string `json:"search"`
		SortType  string `json:"sort_type"`
		Ascend    bool   `json:"ascend_sort"`
	}{PageIndex: 0, PageSize: 20, Ascend: true}
	CheckHandler(!c.ParseBody(&request), message.JsonParseError)
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)

	option := db.NewOptions()
	search := make(map[db.OptionKey]interface{})
	searchAll := make(map[db.OptionKey]interface{})
	search[db.OptType] = LogIn
	searchAll[db.OptType] = []string{LogOut, LogIn}
	var userList []*model.User

	sc := service.GetContainerInstance()
	if userInfo.User.Guid == model.AdminGuid {
		var userIDList []string
		opt := db.NewOptions()
		userList, _, _ = sc.DB.LoadUser(context.TODO(), opt)
		for _, v := range userList {
			if ok, err := regexp.MatchString(request.SearchKey, v.Guid); err == nil && ok {
				userIDList = append(userIDList, v.Guid)
			} else if ok, err := regexp.MatchString(request.SearchKey, v.Name); err == nil && ok {
				userIDList = append(userIDList, v.Guid)
			} else if len(request.SearchKey) == 0 {
				userIDList = append(userIDList, v.Guid)
			}
		}
		search[db.OptUserGuid] = userIDList
		searchAll[db.OptUserGuid] = userIDList
	} else {
		userList = append(userList, userInfo.User)
		search[db.OptUserGuid] = userInfo.User.Guid
		searchAll[db.OptUserGuid] = userInfo.User.Guid
	}
	option.TimeKey = db.OptCreateTime
	option.TimeStart = request.TimeStart
	option.TimeEnd = request.TimeEnd
	option.PageSize = request.PageSize
	option.PageIndex = request.PageIndex
	option.Ascend = request.Ascend
	switch request.SortType {
	case "user_id":
		option.Sort = append(option.Sort, db.OptUserGuid)
	case "time":
		option.Sort = append(option.Sort, db.OptCreateTime)
	default:
		option.Sort = append(option.Sort, db.OptCreateTime)
	}

	option.Search = search
	list, count, err := sc.DB.LoadUserOperation(context.TODO(), option)
	CheckHandler(err, message.GetError)

	option.Search = searchAll
	option.PageSize = 0
	listAll, _, err := sc.DB.LoadUserOperation(context.TODO(), option)
	CheckHandler(err, message.GetError)

	type Response struct {
		UserID     string    `json:"user_id"`
		UserName   string    `json:"user_name"`
		Online     bool      `json:"online"`
		IP         string    `json:"ip_address"`
		LoginTime  time.Time `json:"login_time"`
		LogoutTime time.Time `json:"logout_time"`
	}
	userIDMap := make(map[string][]*model.UserOperation)
	for _, v := range listAll {
		userIDMap[v.UserGuid] = append(userIDMap[v.UserGuid], v)
	}
	userIDNameMap := make(map[string]string)
	userOnlineMap := make(map[string]bool)
	for _, v := range userList {
		userIDNameMap[v.Guid] = v.Name
		option := db.NewOptions()
		option.Sort = []db.OptionKey{db.OptCreateTime}
		option.Ascend = false
		option.PageSize = 1
		option.PageIndex = 0
		option.Search[db.OptUserGuid] = v.Guid
		option.Search[db.OptType] = []string{LogOut, LogIn}
		tempList, _, err := sc.DB.LoadUserOperation(context.TODO(), option)
		CheckHandler(err, message.GetError)
		t := new(model.UserOperation)
		for _, v := range tempList {
			if t.CreateTime.Unix() < v.CreateTime.Unix() {
				t = v
			}
		}
		if t.Type == LogIn {
			userOnlineMap[v.Guid] = true
		} else if t.Type == LogOut {
			userOnlineMap[v.Guid] = false
		}
	}
	var responseList []Response
	for _, v := range list {
		tn := time.Now()
		if request.TimeEnd != 0 {
			tn = time.Unix(request.TimeEnd, 0)
		}
		logoutTime := tn
		for _, w := range userIDMap[v.UserGuid] {
			if w.CreateTime.Unix() > v.CreateTime.Unix() && logoutTime.Unix() > w.CreateTime.Unix() {
				logoutTime = w.CreateTime
			}
		}
		response := Response{}
		response.UserID = v.UserGuid
		response.UserName = userIDNameMap[v.UserGuid]
		response.IP = v.IP
		response.Online = userOnlineMap[v.UserGuid]
		response.LoginTime = v.CreateTime
		if logoutTime == tn {
			response.LogoutTime = time.Unix(0, 0)
			response.Online = true
		} else {
			response.LogoutTime = logoutTime
		}
		responseList = append(responseList, response)
	}
	HttpReponseListHandler(c, count, responseList)
}
