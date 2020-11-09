package impl

import (
	"context"
	"encoding/json"
	framework "tumor_server/framework"
	message "tumor_server/message"
	"tumor_server/model"
	service "tumor_server/service"
)

func GetTask(c *framework.Context) {
	defer PanicHandler(c)
	id := c.GetParam("id")
	sc := service.GetContainerInstance()
	task := sc.DB.GetTask(context.TODO(), id)
	HttpReponseHandler(c, task)
}

func EditTask(c *framework.Context) {
	defer PanicHandler(c)
	data := &model.Task{}
	CheckHandler(!c.ParseBody(data), message.JsonParseError)
	sc := service.GetContainerInstance()
	dbData := sc.DB.GetTask(context.TODO(), data.Guid)
	CheckHandler(!c.ParseBody(dbData), message.JsonParseError)
	err := sc.DB.UpdateTask(context.TODO(), dbData)
	CheckHandler(err, message.UpdateError)
	msgByte, _ := json.Marshal(dbData)
	service.PublishToMQ(dbData.InstitutionId, model.MQTask, msgByte)
	HttpReponseHandler(c, nil)
}
