package impl

import (
	"context"
	"strconv"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	service "tumor_server/service"
)

func GetAppConfig(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	lsit, err := sc.DB.LoadAppConfig(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, lsit)
}

func UserList(c *framework.Context) {
	defer PanicHandler(c)
	sc := service.GetContainerInstance()

	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	search := c.GetURLParam("search")
	ascond := c.GetURLParam("ascond")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)

	option := db.NewOptions()
	option.PageSize = pageSize
	option.PageIndex = pageIndex
	if search != "" {
		option.Match[db.OptName] = search
		option.Match[db.OptPhone] = search
		option.Match[db.OptIDCard] = search
	}
	option.Sort = []db.SortOption{{Key: db.OptID, Ascend: ascond == "true"}}
	lsit, _, err := sc.DB.LoadUser(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, lsit)
}

func UserOperation(c *framework.Context) {
	defer PanicHandler(c)
	pageSizeStr := c.GetURLParam("page_size")
	pageIndexStr := c.GetURLParam("page_index")
	pageSize, err := strconv.Atoi(pageSizeStr)
	CheckHandler(err, message.RequestDataError)
	pageIndex, err := strconv.Atoi(pageIndexStr)
	CheckHandler(err, message.RequestDataError)
	search := c.GetURLParam("search")
	ascond := c.GetURLParam("ascond")
	institution := c.GetURLParam("institution")
	user := c.GetURLParam("user")
	userInfo := &service.UserTokenInfo{}
	c.ParseExtra(userInfo)
	sc := service.GetContainerInstance()
	option := db.NewOptions()
	if institution != "" {
		option.EQ[db.OptInstitution] = institution
	}
	if user != "" {
		option.EQ[db.OptUser] = user
	}
	option.Match[db.OptName] = search
	option.PageSize = pageSize
	option.PageIndex = pageIndex
	option.Sort = []db.SortOption{{Key: db.OptID, Ascend: ascond == "true"}}
	listAll, count, err := sc.DB.LoadUserOperation(context.TODO(), option)
	CheckHandler(err, message.GetError)
	HttpReponseListHandler(c, count, listAll)
}
