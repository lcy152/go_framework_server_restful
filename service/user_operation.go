package service

import (
	"context"
	"time"
	"tumor_server/model"
)

func AddUserRecord(guid, opType, ref, ip string) {
	uop := &model.UserOperation{
		Guid:       NewUUID(),
		UserGuid:   guid,
		Type:       opType,
		RefGuid:    ref,
		CreateTime: time.Now(),
		IP:         ip,
	}
	sc := GetContainerInstance()
	sc.DB.AddUserOperation(context.TODO(), uop)
}
