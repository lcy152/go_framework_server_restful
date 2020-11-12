package db

type tableList struct {
	Institution               string
	InstitutionApplication    string
	AddInstitutionApplication string
	UserToInstitution         string
	UserToInstitutionToRole   string
	User                      string
	UserToUser                string
	UserToUserGroup           string
	DipperMssage              string
	GroupChat                 string
	SingleChat                string
	UserApplication           string
	UserGroup                 string
	UserOperation             string
	AppConfig                 string
	Task                      string
	Role                      string
	Authority                 string
	RoleToAuthority           string
}

var table = tableList{
	Institution:               "institution",
	InstitutionApplication:    "institution_application",
	AddInstitutionApplication: "add_institution_application",
	UserToInstitution:         "user_to_institution",
	UserToInstitutionToRole:   "user_to_institution_to_role",
	User:                      "user",
	UserToUser:                "user_to_user",
	UserToUserGroup:           "user_to_user_group",
	DipperMssage:              "dipper_message",
	GroupChat:                 "group_chat",
	SingleChat:                "single_chat",
	UserApplication:           "user_application",
	UserGroup:                 "user_group",
	UserOperation:             "user_operation",
	AppConfig:                 "app_config",
	Task:                      "task",
	Role:                      "role",
	Authority:                 "authority",
	RoleToAuthority:           "role_to_authority",
}
