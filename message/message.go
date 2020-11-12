package message

const (
	// common
	HttpOk                  = "success"
	HttpError               = "error"
	UnknownError            = "unknown error"
	HttpMethodError         = "method error"
	PhoneError              = "phone exist_error"
	PasswordError           = "password error"
	ValidateError           = "validate error"
	OtherPlaceLoginError    = "other place login error"
	TokenGenerateError      = "token generate error"
	JsonParseError          = "json parse error"
	JsonFormatError         = "json format error"
	RequestDataError        = "request data error"
	RequestRepeatError      = "request repeat error"
	RequestIPForbiddenError = "request forbidden error"
	AuthorityError          = "authority error"
	ImageSizeError          = "image size error"
	// common oprete
	AddError     = "add error"
	UpdateError  = "update error"
	GetError     = "get error"
	DeleteError  = "delete error"
	ExpiredError = "expired error"
	GetListError = "get list error"
	// short message
	ShortMessageFrequentlyError = "get message frequently error"
	ShortMessageLimitedError    = "get message limited error"
	ShortMessageValidateError   = "message validate error"
	// user
	FindUserError   = "user not exist"
	RegisterError   = "register error"
	UpdateUserError = "edit user error"
	DeleteUserError = "delete user error"
	UserExistError  = "user exist"
)
