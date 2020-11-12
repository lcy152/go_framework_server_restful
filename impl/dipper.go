package impl

import (
	"context"
	"tumor_server/db"
	framework "tumor_server/framework"
	message "tumor_server/message"
	service "tumor_server/service"
)

func LoadDipperMessage(c *framework.Context) {
	defer PanicHandler(c)
	var data = struct {
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
		Ascend    bool   `json:"ascend_sort"`
		Search    string `json:"search"`
	}{
		PageSize: 20,
	}
	userInfo := GetContextUserInfo(c)
	sc := service.GetContainerInstance()
	opt := db.NewOptions()
	opt.Match[db.OptSender] = userInfo.User.ID
	opt.Match[db.OptReceiver] = userInfo.User.ID
	if data.Search != "" {
		opt.EQ[db.OptData] = data.Search
	}
	// opt.Sort = append(opt.Sort, db.OptCreateTime)
	msgList, err := sc.DB.LoadDipperMessage(context.TODO(), opt)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, msgList)
}
