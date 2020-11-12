package service

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"tumor_server/db"
	"tumor_server/model"
)

func HandleMessage(data []byte, exchange string) bool {
	request := &model.MqMessage{}
	err := json.Unmarshal(data, request)
	if err != nil {
		log.Printf("Received a message: %s", request.Flag)
		log.Println(err)
		return false
	}
	log.Printf("exchange: %s, type: %s， data length: %d", exchange, request.Flag, len(request.Data))
	switch request.Flag {
	case model.MQUserValidation:
		return ValidateDipperUserConsume(exchange, []byte(request.Data))
	case model.MQTask:
		return MQTaskConsume(exchange, []byte(request.Data))
	case model.MQMessage:
		return MQMessageConsume(exchange, []byte(request.Data))
	}
	return true
}

func ValidateDipperUserConsume(institutionId string, data []byte) bool {
	sc := GetContainerInstance()
	res := &model.ValidateDipperUser{}
	err := json.Unmarshal(data, res)
	if err != nil {
		log.Printf("Received a message: %s", institutionId)
		log.Println(err)
		return false
	}
	opt := db.NewOptions()
	opt.EQ[db.OptInstitution] = res.Institution
	opt.EQ[db.OptUser] = res.UserGuid
	refIns, _ := sc.DB.LoadUserToInstitution(context.TODO(), opt)
	if refIns == nil {
		return false
	}
	if res.Pass {
		for _, v := range refIns {
			v.DipperUser = res.DipperUser
			v.DipperPassword = res.Password
			err := sc.DB.UpdateUserToInstitution(context.TODO(), v)
			if err != nil {
				continue
			}
		}
	} else {
		for _, v := range refIns {
			v.DipperUser = ""
			v.DipperPassword = ""
			err := sc.DB.UpdateUserToInstitution(context.TODO(), v)
			if err != nil {
				continue
			}
		}
	}
	return true
}

func MQTaskConsume(institutionId string, data []byte) bool {
	sc := GetContainerInstance()
	res := &model.Task{}
	err := json.Unmarshal(data, res)
	if err != nil {
		log.Printf("Received a message: %s", string(data))
		log.Println(err)
		return true
	}
	res.ID = NewUUID()

	refIns, _ := sc.DB.GetUserToInstitutionDipperUser(context.TODO(), res.Institution, res.ExecuteUser.String())
	if refIns == nil {
		log.Println(err)
		return false
	}
	user, _ := sc.DB.GetUser(context.TODO(), refIns.User)
	if user == nil {
		return false
	}
	res.InstitutionName = refIns.InstitutionName
	res.ExecuteUser = user.ID
	res.ExecuteUserName = user.Name
	res.ExecuteUserType = refIns.Type

	refIns2, _ := sc.DB.GetUserToInstitutionDipperUser(context.TODO(), res.Institution, res.RefPatientPid)
	if refIns2 != nil {
		res.RefPatient = refIns2.User
		patient, _ := sc.DB.GetUser(context.TODO(), refIns2.User)
		if user != nil {
			res.RefPatientName = patient.Name
		}
	}
	task := sc.DB.GetTaskByOrigin(context.TODO(), res.OriginGuid)
	if task == nil {
		err = sc.DB.AddTask(context.TODO(), res)
		if err != nil {
			log.Printf("AddTask error: %s", err)
			return false
		}
	} else {
		err = sc.DB.UpdateTask(context.TODO(), res)
		if err != nil {
			log.Printf("UpdateTask error: %s", err)
			return false
		}
	}
	return true
}

func MQMessageConsume(institutionId string, data []byte) bool {
	return AddDipperMessage(data)
}

func SendTaskMessage(task *model.Task) {
	taskNameMap := map[string]string{
		model.RegisterTaskType: "登记任务", model.DiagnosisTaskType: "诊断任务", model.LocateApplyTaskType: "定位申请任务",
		model.LocateAppointmentTaskType: "定位预约任务", model.LocateScheduleTaskType: "定位排程任务", model.LocatePatientTaskType: "患者定位任务",
		model.PrescriptionTaskType: "勾画任务", model.PlanDesignTaskType: "计划设计任务", model.PlanApproveTaskType: "计划批准任务",
		model.TreatmentScheduleTaskType: "治疗排程任务", model.TreatmentRecordTaskType: "治疗记录任务",
		model.TreatmentReviewTaskType: "治疗评估任务", model.DischargeRegisterTaskType: "出院登记任务", model.FollowupPlanTaskType: "随访计划任务",
		model.FollowupAppointmentTaskType: "随访预约任务", model.FollowupRecordTaskType: "随访记录任务",
	}
	statusMap := map[string]string{
		model.CONST_TASK_START_CODE: "已创建", model.CONST_TASK_PAUSE_CODE: "已完成", model.CONST_TASK_FINISH_CODE: "已暂停", model.CONST_TASK_REVERT_CODE: "已退回",
	}
	msg := new(model.DipperMessage)
	msg.ID = NewUUID()
	msg.Flag = "task"
	msg.Type = task.TaskType
	// msg.Institution = task.Institution
	msg.InstitutionName = task.InstitutionName
	// msg.RefPatient = task.RefPatientGuid
	msg.RefPatientName = task.RefPatientName
	msg.CreateTime = time.Now()
	msg.Msg = "任务提醒: " + msg.RefPatientName + " 的 " + taskNameMap[task.TaskType] + " " + statusMap[task.TaskState] + "。" + time.Now().Format("2006-01-02 15:04:05")
	msgByte, _ := json.Marshal(msg)
	AddDipperMessage(msgByte)
	return
}

func AddDipperMessage(data []byte) bool {
	sc := GetContainerInstance()
	request := &model.DipperMessage{}
	err := json.Unmarshal(data, request)
	if err != nil {
		log.Printf("Received a message: %s", string(data))
		log.Println(err)
		return true
	}
	err = sc.DB.AddDipperMessage(context.TODO(), request)
	if err != nil {
		log.Printf("AddDipperMessage error, message: %s", string(data))
		return false
	}
	// SendCenterMessage(request.ReceiverGuid, string(data))
	return true
}
