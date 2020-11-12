package main

import (
	framework "tumor_server/framework"
	impl "tumor_server/impl"
)

func NewServer() *framework.Engine {
	s := framework.NewEngine()
	// http
	s.AddMiddleware("", impl.AllRouteMiddleware)
	// static
	s.POST("/app_config", impl.GetAppConfig)
	// file
	s.Static("./dist")
	s.FsFile("/fs", "./fs")

	// v1
	s.AddMiddleware("/v1", impl.V1RouteMiddleware)
	// auth
	s.POST("/v1/login", impl.Login)
	s.POST("/v1/logout", impl.LogOut)
	// short_message
	s.POST("/v1/short_message_login", impl.LoginCode)
	s.POST("/v1/short_message_register", impl.Register)
	s.GET("/v1/short_message/:phone", impl.GetShortMessageCode)

	// admin
	s.AddMiddleware("/v1/admin", impl.V1AdminMiddleware)
	s.GET("/v1/admin/user_list", impl.UserList)
	s.GET("/v1/admin/user_operation", impl.UserOperation)
	s.PUT("/v1/admin/user", impl.AddUser)

	// user
	s.AddMiddleware("/v1/auth", impl.V1AuthMiddleware)
	s.GET("/v1/auth/user", impl.GetUser)
	s.POST("/v1/auth/user", impl.EditUser)

	// user_detail
	s.POST("/v1/auth/user_detail/password", impl.EditUserPassword)
	s.POST("/v1/auth/user_detail/phone", impl.EditUserPhone)
	s.POST("/v1/auth/user_detail/institution/:institution_id", impl.ChangeCurrentInstitution)
	s.GET("/v1/auth/user_detail/institution/:institution_id", impl.GetUserDetailInstitution)
	s.GET("/v1/auth/user_detail/friend_list", impl.UserFriendList)
	s.GET("/v1/auth/user_detail/institution_list", impl.UserInstitutionList)
	s.GET("/v1/auth/user_detail/role_list/:user_id/:institution_id", impl.UserRoleList)
	s.GET("/v1/auth/user_detail/user_group_list", impl.UserUserGroupList)

	// friend
	s.PUT("/v1/auth/friend", impl.AddFriend)
	s.DELETE("/v1/auth/friend/:id", impl.DeleteFriend)
	// friend_detail
	s.GET("/v1/auth/friend_detail/application", impl.GetFriendApplicationList)
	s.POST("/v1/auth/friend_detail/application/:id", impl.ApproveFriend)
	s.DELETE("/v1/auth/friend_detail/application/:id", impl.RejectFriend)

	// institution
	s.GET("/v1/auth/institution/:id", impl.GetInstitution)
	s.PUT("/v1/auth/institution", impl.AddInstitution)
	s.POST("/v1/auth/institution", impl.EditInstitution)
	s.DELETE("/v1/auth/institution/:id", impl.DeleteInstitution)

	// institution detail
	s.GET("/v1/auth/institution_list", impl.LoadInstitution)
	s.GET("/v1/auth/institution_detail/user_list/:institution_id", impl.InstitutionUserList)
	s.GET("/v1/auth/institution_detail/role_list/:institution_id", impl.InstitutionRoleList)

	// institution apply
	s.GET("/v1/auth/institution_detail/application/:institution_id/:state", impl.InstitutionApply)
	s.PUT("/v1/auth/institution_detail/application", impl.ApplyInstitution)
	s.POST("/v1/auth/institution_detail/application/:id", impl.ApproveInstitution)
	s.DELETE("/v1/auth/institution_detail/application/:id", impl.RejectInstitution)

	// chat
	s.GET("/v1/auth/chat_message_history", impl.GetUserChatHistory)
	s.GET("/v1/auth/single_chat_history", impl.LoadSingleChatMessage)
	s.GET("/v1/auth/group_chat_history", impl.LoadGroupChatMessage)
	s.PUT("/v1/auth/single_chat", impl.AddSingleChatMessage)
	s.PUT("/v1/auth/group_chat", impl.AddGroupChatMessage)
	s.DELETE("/v1/auth/single_chat/:id", impl.DeleteSingleChatMessage)
	s.DELETE("/v1/auth/group_chat/:id", impl.DeleteGroupChatMessage)

	// user_group
	s.GET("/v1/auth/user_group/:id", impl.GetUserGroup)
	s.PUT("/v1/auth/user_group", impl.AddUserGroup)
	s.POST("/v1/auth/user_group", impl.EditUserGroup)
	s.DELETE("/v1/auth/user_group/:id", impl.DeleteUserGroup)
	s.GET("/v1/auth/user_group_member/:id", impl.GetUserGroupUser)
	s.PUT("/v1/auth/user_group_member", impl.AddUserGroupUser)
	s.DELETE("/v1/auth/user_group_member/:id", impl.DeleteUserGroupUser)

	// role
	s.GET("/v1/auth/role/:id", impl.GetRole)
	s.DELETE("/v1/auth/role/:id", impl.GetRole)
	s.PUT("/v1/auth/role", impl.AddRole)
	s.POST("/v1/auth/role", impl.EditRole)
	s.GET("/v1/auth/role_detail/authority_list/:role_id", impl.RoleAuthorityList)

	// task
	s.POST("/v1/auth/task", impl.EditTask)
	s.GET("/v1/auth/task/:id", impl.GetTask)
	// search
	s.GET("/v1/auth/search/user", impl.SearchUser)

	// websocket
	s.WS("/v1/ws/message/:token", impl.WSMessage)

	// dipper
	s.AddMiddleware("/v1/dipper", impl.V1DipperMiddleware)
	s.GET("/v1/dipper/message", impl.LoadDipperMessage)

	return s
}
