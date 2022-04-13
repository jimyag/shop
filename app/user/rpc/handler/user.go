package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/user/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

//
// UserServer
//  @Description: user server 的实例
//
type UserServer struct {
	model.Store
	options *password.Options
}

//
// NewUserServer
//  @Description: 使用 store password 的 options 创建 userServer
//  @param store 存储的方式
//  @param options 密码相关的选项
//  @return *UserServer
//
func NewUserServer(store model.Store, options *password.Options) *UserServer {
	return &UserServer{
		Store:   store,
		options: options,
	}
}

//
// userModel2UserInfoResponse
//  @Description: 将user的model转换为user的响应
//  @param user
//  @return *proto.UserInfoResponse
//
func userModel2UserInfoResponse(user model.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:        int32(user.ID),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
		Email:     user.Email,
		Password:  user.Password,
		Nickname:  user.Nickname,
		Gender:    user.Gender,
		Role:      int32(user.Role),
	}
}

//
// GetUserList
//  @Description: 获得用户信息列表
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.UserListResponse
//  @return error
//
func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageIngo) (*proto.UserListResponse, error) {
	arg := model.ListUsersParams{
		Limit:  int32(req.PageSize),
		Offset: int32((req.PageNum - 1) * req.PageSize),
	}

	users, err := u.Store.ListUsers(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获得用户列表信息失败")
	}

	userResponseList := make([]*proto.UserInfoResponse, len(users))
	for i, user := range users {
		userResponseList[i] = userModel2UserInfoResponse(user)
	}
	rsp := proto.UserListResponse{
		Total: int32(len(users)),
		Data:  userResponseList,
	}

	return &rsp, nil
}

//
// GetUserByEmail
//  @Description: 通过 Email 获得用户信息
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.UserInfoResponse
//  @return error
//
func (u *UserServer) GetUserByEmail(ctx context.Context, req *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	// 开始追踪
	// 从context总拿到parentSpan
	parentSpan := opentracing.SpanFromContext(ctx)
	// 生成一个span并设置它的父亲
	getUserInfoByEmailSpan := opentracing.GlobalTracer().
		StartSpan(
			"get user info by email form database",
			opentracing.ChildOf(parentSpan.Context()),
		)
	user, err := u.Store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "通过 Email 获得用户信息失败")
	}
	getUserInfoByEmailSpan.Finish()

	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}

//
// GetUserById
//  @Description: 通过 uid 获得用户信息
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.UserInfoResponse
//  @return error
//
func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	user, err := u.Store.GetUserById(ctx, int64(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "通过 ID 获得用户信息失败")
	}
	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}

//
// CreateUser
//  @Description: 创建用户 创建的时候只需要穿原始的密码即可，这里会进行加密
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.UserInfoResponse
//  @return error
//
func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserInfoResponse, error) {
	salt, pwd := password.Encode(req.Password, u.options)
	createUserParams := model.CreateUserParams{
		Email:    req.GetEmail(),
		Password: fmt.Sprintf("$%s$%s$%s", "pbkdf2-sha512", salt, pwd),
		Nickname: req.GetNickname(),
		Gender:   req.GetGender(),
		Role:     int64(req.GetRole()),
	}

	user, err := u.Store.GetUserByEmail(ctx, createUserParams.Email)
	if err == nil {
		return &proto.UserInfoResponse{}, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "系统错误")
	}

	err = u.Store.ExecTx(ctx, func(queries *model.Queries) error {
		var err error
		user, err = queries.CreateUser(ctx, createUserParams)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// todo 处理Tx的错误
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "系统错误")
	}
	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}

//
// UpdateUser
//  @Description: 更新用户的信息，如果有的话才会更新 字段为空的话保持原来的
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.UserInfoResponse
//  @return error
//
func (u *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserInfoResponse, error) {
	arg := model.UpdateUserParams{
		UpdatedAt: time.Now(),
		Nickname:  req.GetNickname(),
		Gender:    req.GetGender(),
		Role:      int64(req.GetRole()),
		ID:        int64(req.Id),
	}

	user, err := u.Store.GetUserById(ctx, arg.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return &proto.UserInfoResponse{}, status.Errorf(codes.NotFound, "用户不存在")
	} else if err != nil {
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "系统错误")
	}
	// 如果要更新密码 重新设置
	if req.GetPassword() != "" {
		salt, pwd := password.Encode(req.Password, u.options)
		arg.Password = fmt.Sprintf("$%s$%s$%s", "pbkdf2-sha512", salt, pwd)
	}

	err = u.Store.ExecTx(ctx, func(queries *model.Queries) error {
		if req.Nickname == "" {
			arg.Nickname = user.Nickname
		}
		if req.Gender == "" {
			arg.Gender = user.Gender
		}
		if req.Role == 0 {
			arg.Gender = user.Gender
		}
		if req.Password == "" {
			arg.Password = user.Password
		}

		user, err = queries.UpdateUser(ctx, arg)
		return err
	})
	if err != nil {
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "系统错误")
	}
	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}

//
// CheckPassword
//  @Description: 检查用户的密码是否合规
//  @receiver u
//  @param ctx
//  @param req
//  @return *proto.CheckPasswordResponse
//  @return error
//
func (u *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckPasswordResponse, error) {
	// 开始追踪
	// 从context总拿到parentSpan
	parentSpan := opentracing.SpanFromContext(ctx)
	// 生成一个span并设置它的父亲
	checkUserPassword := opentracing.GlobalTracer().
		StartSpan(
			"check user password",
			opentracing.ChildOf(parentSpan.Context()),
		)

	encryptedPasswordInfo := strings.Split(req.GetEncryptedPassword(), "$")
	check := password.Verify(req.Password, encryptedPasswordInfo[2], encryptedPasswordInfo[3], u.options)

	checkUserPassword.Finish()
	return &proto.CheckPasswordResponse{Success: check}, nil
}
