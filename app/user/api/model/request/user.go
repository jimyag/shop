package request

//
// PasswordLoginForm
//  @Description: 使用邮箱和密码登录的参数
//
type PasswordLoginForm struct {
	Email    string `json:"email" validate:"required,email" label:"邮件"`
	Password string `json:"password" validate:"required,min=6,max=20" label:"您的密码"`
}

//
// CreateUser
//  @Description: 创建用户的参数
//
type CreateUser struct {
	Email      string `json:"email" validate:"required,email" label:"邮件"`
	Password   string `json:"password" validate:"required,min=6,max=20" label:"您的密码"`
	RePassword string `json:"re_password" validate:"required,eqfield=Password" label:"确认密码"`
	Nickname   string `json:"nickname" validate:"required,min=6,max=20" label:"昵称"`
	Gender     string `json:"gender" validate:"required,oneof=male femal" label:"性别"`
	Role       int32  `json:"role" validate:"required,oneof=1 2" label:"权限"`
	AuthCode   int    `json:"auth_code" validate:"required,min=10000,max=99999" label:"验证码"`
}

//
// CreateUserEmail
//  @Description: 创建用户时发送邮件的参数
//
type CreateUserEmail struct {
	Email string `json:"email" validate:"required,email" label:"邮件"`
}

//
// GetUserByEmail
//  @Description: 通过邮箱获得用户信息的参数
//
type GetUserByEmail struct {
	Email string `json:"email" validate:"required,email" label:"邮箱"`
}

//
// GetUserByID
//  @Description: 通过uid获得用户信息的参数
//
type GetUserByID struct {
	ID uint32 `json:"ID" validate:"required,min=1"`
}

//
// UpdateUserWithoutPwd
//  @Description: 更新用户的nickname 和gender
//
type UpdateUserWithoutPwd struct {
	Id       int32  `json:"id" validate:"required,min=1"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
}

//
// ChangePassword
//  @Description: 修改用的密码
//
type ChangePassword struct {
	Id         int32  `json:"id" validate:"required,min=1"`
	Password   string `json:"password" validate:"required,min=6,max=20" label:"您的密码"`
	RePassword string `json:"re_password" validate:"required,eqfield=Password" label:"确认密码"`
}

//
// ChangeRole
//  @Description: 修改用户的权限
//
type ChangeRole struct {
	Id   int32 `json:"id" validate:"required,min=1"`
	Role int32 `json:"role" validate:"required,oneof=1 2" label:"权限"`
}
