package api

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/global"
	"github.com/jimyag/shop/app/user/api/model/request"
	"github.com/jimyag/shop/app/user/api/model/response"
	"github.com/jimyag/shop/app/user/api/tools/auth_code"
	"github.com/jimyag/shop/app/user/api/tools/email"
	"github.com/jimyag/shop/app/user/api/tools/paseto"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/handle_grpc_error"
	"github.com/jimyag/shop/common/utils/validate"
)

//
// GetUserList
//  @Description: 获得用户的列表
//  @param ctx
//
func GetUserList(ctx *gin.Context) {
	// 处理请求的参数
	userListArg := request.GetUserList{}
	_ = ctx.ShouldBindQuery(&userListArg)
	msg, err := validate.Validate(userListArg, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 查询用户
	rsp, err := global.UserSrvClient.GetUserList(ctx, &proto.PageIngo{
		PageNum:  uint32(userListArg.PageNum),
		PageSize: uint32(userListArg.PageSize),
	})
	if err != nil {
		global.Logger.Error("获得用户列表失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成结果
	result := make([]interface{}, 0)
	for _, datum := range rsp.Data {
		data := make(map[string]interface{})
		data["id"] = datum.Id
		data["nickname"] = datum.Nickname
		data["gender"] = datum.Gender
		data["email"] = datum.Email
		data["role"] = datum.Role
		result = append(result, data)
	}
	res := make(map[string]interface{})
	res["total"] = rsp.Total
	res["data"] = result
	model.OkWithData(res, ctx)
}

//
// PasswordLogin
//  @Description: 使用邮箱和密码登录
//  @param ctx
//
func PasswordLogin(ctx *gin.Context) {
	// 处理参数
	passwordLoginForm := request.PasswordLoginForm{}
	_ = ctx.ShouldBindJSON(&passwordLoginForm)
	msg, err := validate.Validate(passwordLoginForm, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 使用email查询用户
	user, err := global.UserSrvClient.GetUserByEmail(ctx, &proto.EmailRequest{Email: passwordLoginForm.Email})
	if err != nil {
		global.Logger.Info("查找用户错误", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("用户不存在", ctx)
		return
	}

	// 检查用户的密码
	checkP := proto.PasswordCheckInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: user.GetPassword(),
	}
	password, err := global.UserSrvClient.CheckPassword(ctx, &checkP)
	if err != nil {
		global.Logger.Info("用户登录错误", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("登录失败", ctx)
		return
	}
	if !password.GetSuccess() {
		model.FailWithMsg("邮箱或密码错误", ctx)
		return
	}

	// 生成token
	payload, _ := paseto.NewPayload(user.Id, user.Role)
	token, err := global.PasetoMaker.CreateToken(payload)
	if err != nil {
		global.Logger.Info("创建Token失败", zap.Error(err))
		model.FailWithMsg("登录失败", ctx)
		return
	}

	// 响应
	res := make(map[string]string)
	res["token"] = token
	model.OkWithDataMsg(res, "登录成功", ctx)
}

//
// CreateUserEmail
//  @Description: 创建用户发送验证码
//  @param ctx
//
func CreateUserEmail(ctx *gin.Context) {
	// 校验参数
	createUserEmail := request.CreateUserEmail{}
	_ = ctx.ShouldBindJSON(&createUserEmail)
	msg, err := validate.Validate(createUserEmail, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 判断用户是否已经发送过验证码了
	i, err := global.Redis.CreateUser.Exists(context.Background(), createUserEmail.Email).Result()
	if i == 1 {
		model.FailWithMsg("验证码已经发送，请稍后", ctx)
		return
	}

	// 没有发送
	code := auth_code.Code()
	_, err = global.Redis.CreateUser.SetEX(context.Background(),
		createUserEmail.Email,
		code,
		time.Duration(global.RemoteConfig.Timeout.CreateUserEmail)*time.Minute,
	).Result()
	if err != nil {
		global.Logger.Info("注册验证码写入redis失败", zap.Error(err))
		model.FailWithMsg("验证码发送失败，请稍后重试", ctx)
		return
	}

	// 邮件的内容
	header := `<h1>【管理系统】 您正在 注册账号  。您的验证码是：`
	header += strconv.Itoa(code) + ` ，请勿向他人泄露，转发可能导致账号被盗。如非本人操作，可忽略本消息。5分钟内有效</h1>`
	err = email.Send([]string{createUserEmail.Email}, "用户注册", header)
	if err != nil {
		global.Logger.Info("邮件发送失败", zap.Error(err))
		model.FailWithMsg("验证码发送失败,请稍后重试", ctx)
		return
	}

	model.OkWithMsg("验证码发送成功，请注意查收", ctx)
}

//
// CreateUser
//  @Description: 创建用户
//  @param ctx
//
func CreateUser(ctx *gin.Context) {
	// 校验参数
	createUserParams := request.CreateUser{}
	_ = ctx.ShouldBindJSON(&createUserParams)
	msg, err := validate.Validate(createUserParams, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 是否发送过验证码
	i, _ := global.Redis.CreateUser.Exists(ctx, createUserParams.Email).Result()
	if i == 0 {
		model.FailWithMsg("请先发送验证码", ctx)
		return
	}

	// 有 判断和发送来的是否一样
	i, _ = global.Redis.CreateUser.Get(ctx, createUserParams.Email).Int64()
	if i != int64(createUserParams.AuthCode) {
		model.FailWithMsg("验证码错误", ctx)
		return
	}

	// 创建用户
	// 这里只需要传没经过加密密码的就可以
	createUserRequest := proto.CreateUserRequest{
		Email:    createUserParams.Email,
		Password: createUserParams.RePassword,
		Nickname: createUserParams.Nickname,
		Gender:   createUserParams.Gender,
		Role:     createUserParams.Role,
	}
	_, err = global.UserSrvClient.CreateUser(ctx, &createUserRequest)
	if err != nil {
		global.Logger.Info("创建用户失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("创建用户失败", ctx)
		return
	}

	// 创建用户之后删除验证码
	global.Redis.CreateUser.Del(ctx, createUserParams.Email)
	model.OkWithMsg("创建用户成功", ctx)
}

//
// GetUserByEmail
//  @Description: 使用邮箱获得用户信息
//  @param ctx
//
func GetUserByEmail(ctx *gin.Context) {
	// 校验参数
	userByEmailParam := request.GetUserByEmail{}
	_ = ctx.ShouldBindJSON(&userByEmailParam)
	msg, err := validate.Validate(userByEmailParam, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 获得用户信息
	user, err := global.UserSrvClient.GetUserByEmail(ctx, &proto.EmailRequest{Email: userByEmailParam.Email})
	if err != nil {
		global.Logger.Info("查找用户", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	//
	rsp := response.GetUserInfoResponse{}
	err = copier.Copy(&rsp, &user)
	if err != nil {
		global.Logger.Error("从user -> rsp_user出错", zap.Error(err))
		model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}
	model.OkWithDataMsg(rsp, "成功", ctx)
}

//
// GetUserByID
//  @Description: 通过uid获得用户信息
//  @param ctx
//
func GetUserByID(ctx *gin.Context) {
	// 校验参数
	userByIDParam := request.GetUserByID{}
	_ = ctx.ShouldBindJSON(&userByIDParam)
	msg, err := validate.Validate(userByIDParam, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 获得用户信息
	user, err := global.UserSrvClient.GetUserById(ctx, &proto.IdRequest{Id: userByIDParam.ID})
	if err != nil {
		global.Logger.Error("使用ID查询用户失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}
	// 处理结果
	rsp := response.GetUserInfoResponse{}
	err = copier.Copy(&rsp, &user)
	if err != nil {
		global.Logger.Error("从user -> rsp_user出错", zap.Error(err))
		model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	model.OkWithDataMsg(rsp, "成功", ctx)
}

//
// UpdateUserWithOutPassword
//  @Description: 更新用户的nickname 和gender
//  @param ctx
//
func UpdateUserWithOutPassword(ctx *gin.Context) {
	// 参数校验
	updateUserParams := request.UpdateUserWithoutPwd{}
	_ = ctx.ShouldBindJSON(&updateUserParams)
	msg, err := validate.Validate(updateUserParams, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 更新用户信息
	arg := proto.UpdateUserRequest{}
	err = copier.Copy(&arg, &updateUserParams)
	if err != nil {
		global.Logger.Error("从update_user_param -> update_user_request出错", zap.Error(err))
		model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	user, err := global.UserSrvClient.UpdateUser(ctx, &arg)
	if err != nil {
		global.Logger.Info("更新用户信息出错", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	model.OkWithData(user, ctx)
}

//
// ChangePassword
//  @Description: 修改用户密码
//  @param ctx
//
func ChangePassword(ctx *gin.Context) {
	// 校验参数
	changePassword := request.ChangePassword{}
	_ = ctx.ShouldBindJSON(&changePassword)
	msg, err := validate.Validate(changePassword, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	arg := proto.UpdateUserRequest{}
	err = copier.Copy(&arg, &changePassword)
	if err != nil {
		global.Logger.Error("从changePassword_param -> update_user_request出错", zap.Error(err))
		model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	// 更新用户密码
	user, err := global.UserSrvClient.UpdateUser(ctx, &arg)
	if err != nil {
		global.Logger.Error("更新用户密码出错", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}

	model.OkWithData(user, ctx)
}

//
// ChangeRole
//  @Description: 修改用户权限
//  @param ctx
//
func ChangeRole(ctx *gin.Context) {
	// 参数校验
	changeRole := request.ChangeRole{}
	_ = ctx.ShouldBindJSON(&changeRole)
	msg, err := validate.Validate(changeRole, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	arg := proto.UpdateUserRequest{}
	err = copier.Copy(&arg, &changeRole)
	if err != nil {
		global.Logger.Error("从changeRole_param -> update_user_request出错", zap.Error(err))
		model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}
	user, err := global.UserSrvClient.UpdateUser(ctx, &arg)
	if err != nil {
		global.Logger.Error("更新用户权限", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		//model.FailWithMsg("系统错误，请稍后重试", ctx)
		return
	}
	model.OkWithData(user, ctx)
}
