package db

type OptionKey uint64

const (
	OptInvalid OptionKey = iota
	// base
	OptGuid
	OptName
	OptInstitutionId
	OptInstitutionList
	OptInstitutionName
	OptDipperUser
	OptCreateTime
	OptLastModTime
	OptType
	OptFlag
	OptStatus
	OptDate
	OptPhone
	OptSex
	OptBirthDate
	OptKeyCode
	OptToken
	OptAddress
	OptCreator
	OptCode
	OptIDCard
	// user
	OptApprove
	OptDisable
	OptRole
	OptGroupIdList
	OptGroupId
	OptData
	OptUserID
	OptCurrent
	// group
	OptManager
	OptMember
	// patient
	OptUserGuid
	OptPatientName
	OptPatientGuid
	OptPatientPid
	// report patient
	OptReportTime
	OptBeginTime
	OptIndexNumber
	// study
	OptStudyGuid
	// plan
	OptFxNumber
	OptFxCount
	OptPlanGuid
	// task
	OptTaskGuid
	OptTaskState
	OptTaskType
	// tumor
	OptSenderGuid
	OptReceiverGuid
)
