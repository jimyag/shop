package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/user/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

type UserServer struct {
	model.Store
	options *password.Options
}

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

func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageIngo) (*proto.UserListResponse, error) {
	arg := model.ListUsersParams{
		Limit:  int32(req.PageSize),
		Offset: int32((req.PageNum - 1) * req.PageSize),
	}

	users, err := u.Store.ListUsers(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获得用户列表信息失败")
	}
	// todo 获得total
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
func (u *UserServer) GetUserByEmail(ctx context.Context, req *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	user, err := u.Store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "通过 Email 获得用户信息失败")
	}
	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}
func (u *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	user, err := u.Store.GetUserById(ctx, int64(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "通过 ID 获得用户信息失败")
	}
	rsp := userModel2UserInfoResponse(user)
	return rsp, nil
}

// CreateUser 密码这里处理
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
func (u *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckPasswordResponse, error) {
	encryptedPasswordInfo := strings.Split(req.GetEncryptedPassword(), "$")
	check := password.Verify(req.Password, encryptedPasswordInfo[2], encryptedPasswordInfo[3], u.options)
	return &proto.CheckPasswordResponse{Success: check}, nil
}
