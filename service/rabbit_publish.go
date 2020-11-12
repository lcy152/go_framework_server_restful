package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"tumor_server/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PublishToMQ(institutionId string, flag string, data []byte) {
	sc := GetContainerInstance()
	insID, _ := primitive.ObjectIDFromHex(institutionId)
	ins, _ := sc.DB.GetInstitution(context.TODO(), insID)
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
		sc.RabbitMQ.Publish(ins.ID.String(), string(msgByte))
	}
}

func ValidateDipperUserPublish(du *model.ValidateDipperUser) error {
	if du.Institution.String() == "" {
		log.Println("empty institution")
		return errors.New("empty institutionId")
	}
	msgByte, err := json.Marshal(du)
	if err != nil {
		log.Println(err)
		return err
	}
	PublishToMQ(du.Institution.String(), model.MQUserValidation, msgByte)
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
	res.ID = NewUUID()

	refIns, _ := sc.DB.GetUserToInstitutionUserType(context.TODO(), res.Institution, res.ExecuteUser, model.UserToInstitutionWorker)
	if refIns == nil {
		log.Println(err)
		return false
	}
	_, err = sc.DB.GetUser(context.TODO(), refIns.User)
	if err != nil {
		return false
	}
	// res.ExecuteUser = refIns.DipperUser

	refIns2, _ := sc.DB.GetUserToInstitutionUserType(context.TODO(), res.Institution, res.RefPatient, model.UserToInstitutionWorker)
	if refIns2 == nil {
		res.RefPatient = refIns2.User
	}
	msgByte, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return false
	}
	PublishToMQ(res.Institution.String(), model.MQTask, msgByte)
	return true
}
