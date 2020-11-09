package message

const (
	// common
	HttpOk                  = "success"
	HttpError               = "error"
	HttpMethodError         = "method_error"
	PhoneError              = "phone_exist_errort"
	PasswordError           = "password_errort"
	ValidateError           = "validate_error"
	OtherPlaceLoginError    = "other_place_login_error"
	TokenGenerateError      = "token_generate_error"
	JsonParseError          = "json_parse_error"
	JsonFormatError         = "json_format_error"
	RequestDataError        = "request_data_error"
	RequestRepeatError      = "request_repeat_error"
	RequestIPForbiddenError = "request_forbidden_error"
	AuthorityError          = "authority_error"
	ImageSizeError          = "image_size_error"
	// common oprete
	AddError     = "add_error"
	UpdateError  = "update_error"
	GetError     = "get_error"
	DeleteError  = "delete_error"
	ExpiredError = "expired_error"
	// short message
	ShortMessageFrequentlyError = "get_message_frequently_error"
	ShortMessageLimitedError    = "get_message_limited_error"
	ShortMessageValidateError   = "message_validate_error"
	// user
	FindUserError   = "user_not_exist"
	RegisterError   = "register_error"
	UpdateUserError = "edit_user_error"
	DeleteUserError = "delete_user_error"
	UserExistError  = "user_exist"
)
