package db

type OptionKey uint64

const (
	OptInvalid OptionKey = iota
	// base
	OptID
	OptName
	OptInstitution
	OptInstitutionName
	OptDipperUser
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
	OptCode
	OptIDCard
	// user
	OptDisable
	OptHidden
	OptData
	OptUser
	OptUserName
	OptCurrent
	OptNumber
	// group
	OptManager
	OptMember
	OptGroup
	// patient
	OptFriend
	OptFriendName
	OptPatient
	OptPatientName
	OptPatientPid
	// report patient
	OptReportTime
	OptBeginTime
	OptIndexNumber
	// study
	OptStudy
	// plan
	OptFxNumber
	OptFxCount
	OptPlan
	// task
	OptTask
	OptTaskState
	OptTaskType
	// tumor
	OptSender
	OptReceiver
	// role
	OptRole
	OptAuthority
)
