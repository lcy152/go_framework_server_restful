package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tumor_server/model"
)

func PublishToMQ(institutionId string, flag string, data []byte) {
	sc := GetContainerInstance()
	ins := sc.DB.GetInstitution(context.TODO(), institutionId)
	if ins == nil {
		log.Println("empty institution")
		return
	}
	msg := new(model.MqMessage)
	msg.Flag = flag
	msg.KeyCode = ins.KeyCode
	msg.Data = string(data)
	msgByte, _ := json.Marshal(msg)
	if sc.RabbitMQ != nil {
		sc.RabbitMQ.Publish(ins.Guid, string(msgByte))
	}
}

func ValidateDipperUserPublish(du *model.ValidateDipperUser) error {
	if du.InstitutionId == "" {
		log.Println("empty institution")
		return errors.New("empty institutionId")
	}
	msgByte, err := json.Marshal(du)
	if err != nil {
		log.Println(err)
		return err
	}
	PublishToMQ(du.InstitutionId, model.MQUserValidation, msgByte)
	return nil
}

func DipperTaskPublish(institutionId string, data []byte) bool {
	sc := GetContainerInstance()
	res := &model.Task{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Printf("Received a message: %s", string(data))
		log.Println(err)
		return true
	}
	res.Guid = NewUUID()

	refIns := sc.DB.GetUserRouterByTumorUser(context.TODO(), res.InstitutionId, res.ExecuteUserId)
	if refIns == nil {
		log.Println(err)
		return false
	}
	user := sc.DB.GetUser(context.TODO(), refIns.UserGuid)
	if user == nil {
		return false
	}
	res.ExecuteUserId = refIns.DipperUser

	refIns2 := sc.DB.GetUserRouterByTumorUser(context.TODO(), res.InstitutionId, res.RefPatientGuid)
	if refIns2 == nil {
		res.RefPatientGuid = refIns2.UserGuid
	}
	msgByte, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return false
	}
	PublishToMQ(res.InstitutionId, model.MQTask, msgByte)
	return true
}
