package db

type tableList struct {
	Institution     string
	UserRouter      string
	User            string
	DipperMssage    string
	GroupChat       string
	SingleChat      string
	UserApplication string
	UserGroup       string
	UserOperation   string
	AppConfig       string
	Task            string
}

var table = tableList{
	Institution:     "institution",
	UserRouter:      "user_router",
	User:            "user",
	DipperMssage:    "dipper_message",
	GroupChat:       "group_chat",
	SingleChat:      "single_chat",
	UserApplication: "user_application",
	UserGroup:       "user_group",
	UserOperation:   "user_operation",
	AppConfig:       "app_config",
	Task:            "task",
}
