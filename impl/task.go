package impl

import (
	"context"
	"encoding/json"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTask(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	oid, _ := primitive.ObjectIDFromHex(id)
	sc := service.GetContainerInstance()
	task, err := sc.DB.GetTask(context.TODO(), oid)
	CheckHandler(err, message.GetError)
	HttpReponseHandler(c, task)
}

func EditTask(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Task{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	dbData, err := sc.DB.GetTask(context.TODO(), data.ID)
	CheckHandler(!c.ParseBody(dbData), message.JsonParseError)
	CheckHandler(err, message.GetError)
	err = sc.DB.UpdateTask(context.TODO(), dbData)
	CheckHandler(err, message.UpdateError)
	msgByte, _ := json.Marshal(dbData)
	service.PublishToMQ(dbData.Institution.String(), model.MQTask, msgByte)
	HttpReponseHandler(c, nil)
}
