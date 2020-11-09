package model

import "time"

type Hospital struct {
	Doctor       string    `json:"doctor" bson:"doctor"`
	FxNumber     string    `json:"fx_number" bson:"fx_number"`
	Group        string    `json:"group" bson:"group"`
	ClinicNumber string    `json:"clinic number" bson:"clinic number"` // 门诊号
	Department   string    `json:"department" bson:"department"`       // 科室
	BedNumber    string    `json:"bed_number" bson:"bed_number"`
	Description  string    `json:"description" bson:"description"`
	InTime       time.Time `json:"in_time" bson:"in_time"`
}

type Diagnosis struct {
	Description      string    `json:"description" bson:"description"`
	DiagnosisStaging string    `json:"diagnosis_staging" bson:"diagnosis_staging"` // 诊断分期
	DiagnosisType    string    `json:"diagnosis_type" bson:"diagnosis_type"`       // 诊断类型
	DiagnosisMode    string    `json:"diagnosis_mode" bson:"diagnosis_mode"`       // 诊断方案
	Remark           string    `json:"remark" bson:"remark"`
	DiagnosisDate    time.Time `json:"diagnosis_date" bson:"diagnosis_date"`
}

type LocateApplication struct {
	Description   string   `json:"description" bson:"description"`
	UpperBound    string   `json:"upper_bound" bson:"upper_bound"`
	LowerBound    string   `json:"lower_bound" bson:"lower_bound"`
	Center        string   `json:"center" bson:"center"`
	Mode          string   `json:"mode" bson:"mode"`
	Process       string   `json:"process" bson:"process"`
	LayerSpacing  string   `json:"layer_spacing" bson:"layer_spacing"`
	ScanPosition  string   `json:"scan_position" bson:"scan_position"`
	PillowType    string   `json:"pillow_type" bson:"pillow_type"`
	HandsPosition string   `json:"hands_position" bson:"hands_position"`
	Confirm       bool     `json:"confirm" bson:"confirm"`
	PhotoList     []string `json:"photo_list" bson:"photo_list"`
}

type PlanDesign struct {
	FxCount    float64  `json:"fx_count" bson:"fx_count"`
	TxModality string   `json:"type" bson:"type"`
	Gtv        string   `json:"gtv" bson:"gtv"`
	Ctv        string   `json:"ctv" bson:"ctv"`
	Dose       float64  `json:"dose" bson:"dose"`
	PhotoList  []string `json:"photo_list" bson:"photo_list"`
}

type Treatment struct {
	StartTime   string    `json:"start_time" bson:"start_time"`
	TreatTime   time.Time `json:"treat_time" bson:"treat_time"`
	PhotoList   []string  `json:"photo_list" bson:"photo_list"`
	Description string    `json:"description" bson:"description"`
}

type Task struct {
	Description         string             `json:"description" bson:"description"`
	InstitutionId       string             `json:"institution_id" bson:"institution_id"`
	InstitutionName     string             `json:"institution_name" bson:"institution_name"`
	Guid                string             `json:"guid" bson:"_id"`
	Name                string             `json:"name" bson:"name"`
	TaskType            string             `json:"task_type" bson:"task_type"`
	TaskState           string             `json:"task_state" bson:"task_state"`
	RefPatientPid       string             `json:"ref_patient_pid" bson:"ref_patient_pid"`
	RefPatientGuid      string             `json:"ref_patient_guid" bson:"ref_patient_guid"`
	RefPatientName      string             `json:"ref_patient_name" bson:"ref_patient_name"`
	RefWorkflowGuid     string             `json:"ref_workflow_guid" bson:"ref_workflow_guid"`
	RefDiagnosisGuid    string             `json:"ref_diagnosis_guid" bson:"ref_diagnosis_guid"`
	RefImageGuid        string             `json:"ref_image_guid" bson:"ref_image_guid"`
	RefPlanGuid         string             `json:"ref_plan_guid" bson:"ref_plan_guid"`
	RefFxNumber         int32              `json:"ref_fx_number" bson:"ref_fx_number"`
	RefMachineGuid      string             `json:"ref_machine_guid" bson:"ref_machine_guid"`
	RefMachineName      string             `json:"ref_machine_name" bson:"ref_machine_name"`
	Stage               string             `json:"stage" bson:"stage"`
	NextTaskType        string             `json:"next_task_type" bson:"next_task_type"`
	NextUserName        string             `json:"next_user_name" bson:"next_user_name"`
	NextUserId          string             `json:"next_user_id" bson:"next_user_id"`
	ExecuteUserId       string             `json:"execute_user_id" bson:"execute_user_id"`
	ExecuteUserType     string             `json:"execute_user_type" bson:"execute_user_type"`
	ExecuteUserName     string             `json:"execute_user_name" bson:"execute_user_name"`
	Deadline            time.Time          `json:"deadline" bson:"deadline"`
	Priority            string             `json:"priority" bson:"priority"`
	Comment             string             `json:"comment" bson:"comment"`
	Data                string             `json:"data" bson:"data"`
	RefFollowupPlanGuid string             `json:"ref_followup_plan_guid" bson:"ref_followup_plan_guid"`
	RefStudyGuid        string             `json:"ref_study_guid" bson:"ref_study_guid"`
	RefFlowNumber       int64              `json:"ref_flow_number" bson:"ref_flow_number"`
	RefFlowGuid         string             `json:"ref_flow_guid" bson:"ref_flow_guid"`
	LastTaskGuid        string             `json:"last_task_guid" bson:"last_task_guid"`
	Creator             string             `json:"creator" bson:"creator"`
	LastOperator        string             `json:"last_operator" bson:"last_operator"`
	CreatedTime         time.Time          `json:"created_time" bson:"created_time"`
	LastModTime         time.Time          `json:"last_mod_time" bson:"last_mod_time"`
	TreatmentList       []*Treatment       `json:"treatment_list" bson:"treatment_list"`         // *
	Hospital            *Hospital          `json:"hospital" bson:"hospital"`                     // *
	Diagnosis           *Diagnosis         `json:"diagnosis" bson:"diagnosis"`                   // *
	LocateApplication   *LocateApplication `json:"locate_application" bson:"locate_application"` // *
	PlanDesign          *PlanDesign        `json:"plan_design" bson:"plan_design"`               // *
	PhotoList           []string           `json:"photo_list" bson:"photo_list"`                 // *
	OriginGuid          string             `json:"origin_guid" bson:"origin_guid"`               // *
	OriginTimeStamp     int64              `json:"origin_time_stamp" bson:"origin_time_stamp"`   // *
}

const (
	RegisterTaskType            = "register"
	DiagnosisTaskType           = "diagnosis"
	LocateApplyTaskType         = "locate_apply"
	LocateAppointmentTaskType   = "locate_appointment"
	LocateScheduleTaskType      = "locate_schedule"
	LocatePatientTaskType       = "locate_patient"
	PrescriptionTaskType        = "prescription"
	PlanDesignTaskType          = "plan_design"
	PlanApproveTaskType         = "plan_approve"
	TreatmentScheduleTaskType   = "treatment_schedule"
	TreatmentRecordTaskType     = "treatment_record"
	TreatmentReviewTaskType     = "treatment_review"
	DischargeRegisterTaskType   = "discharge_register"
	FollowupPlanTaskType        = "followup_plan"
	FollowupAppointmentTaskType = "followup_appointment"
	FollowupRecordTaskType      = "followup_record"
)

const (
	CONST_TASK_START_CODE  = "init"
	CONST_TASK_PAUSE_CODE  = "pause"
	CONST_TASK_FINISH_CODE = "finish"
	CONST_TASK_REVERT_CODE = "revert"
)
