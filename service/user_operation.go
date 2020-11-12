package service

import (
	"context"
	"time"
	"tumor_server/db"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUserRecord(guid primitive.ObjectID, opType string, ip string) {
	sc := GetContainerInstance()
	if opType == model.UserOperationLogin {
		option := db.NewOptions()
		option.Sort = []db.SortOption{{Key: db.OptID, Ascend: false}}
		option.PageSize = 1
		option.PageIndex = 0
		option.EQ[db.OptUser] = guid
		option.EQ[db.OptType] = []string{model.UserOperationLogin, model.UserOperationLogout}
		tempList, _, _ := sc.DB.LoadUserOperation(context.TODO(), option)
		if len(tempList) > 0 {
			if tempList[0].Type == model.UserOperationLogin {
				uop := &model.UserOperation{
					ID:         NewUUID(),
					User:       guid,
					Type:       model.UserOperationLogout,
					CreateTime: time.Now(),
					IP:         ip,
				}
				sc.DB.AddUserOperation(context.TODO(), uop)
			}
		}
	}
	uop := &model.UserOperation{
		ID:         NewUUID(),
		User:       guid,
		Type:       opType,
		CreateTime: time.Now(),
		IP:         ip,
	}
	sc.DB.AddUserOperation(context.TODO(), uop)
}
