package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/stretchr/testify/require"

	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/test_util"
)

func createUser(t *testing.T) (*proto.UserInfoResponse, test_util.Password) {
	p := test_util.RandomPassword()
	request := proto.CreateUserRequest{
		Email:    test_util.RandomEmail(),
		Password: p.RawPassword,
		Nickname: test_util.RandomNickName(),
		Gender:   test_util.RandomGender(),
		Role:     0,
	}
	rsp, err := userClient.CreateUser(context.Background(), &request)
	require.NoError(t, err)
	require.NotEmpty(t, rsp)

	require.Equal(t, request.Email, rsp.GetEmail())
	require.Equal(t, request.Nickname, rsp.GetNickname())
	require.Equal(t, request.Gender, rsp.GetGender())
	require.Equal(t, request.Role, rsp.GetRole())
	return rsp, p
}

func TestUserServer_CreateUser(t *testing.T) {
	rsp, p := createUser(t)
	request := proto.CreateUserRequest{
		Email:    rsp.GetEmail(),
		Password: p.RawPassword,
		Nickname: rsp.GetNickname(),
		Gender:   rsp.GetGender(),
		Role:     0,
	}
	newRsp, err := userClient.CreateUser(context.Background(), &request)
	require.Error(t, err)
	require.Empty(t, newRsp)
}

func TestUserServer_GetUserById(t *testing.T) {
	rsp, _ := createUser(t)
	request := proto.IdRequest{Id: uint32(rsp.Id)}
	userRsp, err := userClient.GetUserById(context.Background(), &request)
	require.NoError(t, err)
	require.Equal(t, rsp, userRsp)
}

func TestUserServer_GetUserByEmail(t *testing.T) {
	rsp, _ := createUser(t)
	request := proto.EmailRequest{Email: rsp.Email}
	userRsp, err := userClient.GetUserByEmail(context.Background(), &request)
	require.NoError(t, err)
	require.Equal(t, rsp, userRsp)
}

func TestUserServer_GetUserList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createUser(t)
	}
	request := proto.PageIngo{PageNum: 1, PageSize: 10}
	rsp, err := userClient.GetUserList(context.Background(), &request)
	require.NoError(t, err)
	require.Len(t, rsp.Data, int(request.PageSize))
	require.Equal(t, int(rsp.Total), int(request.PageSize))

	for _, datum := range rsp.Data {
		require.NotEmpty(t, datum)
	}
}

func TestUserServer_UpdateUser(t *testing.T) {
	rsp, _ := createUser(t)
	request := proto.UpdateUserRequest{
		Email:    test_util.RandomEmail(),
		Nickname: test_util.RandomNickName(),
		Gender:   test_util.RandomGender(),
		Role:     0,
		Id:       rsp.GetId(),
	}
	newRsp, err := userClient.UpdateUser(context.Background(), &request)
	require.NoError(t, err)
	require.Equal(t, request.GetNickname(), newRsp.GetNickname())
	require.Equal(t, request.GetGender(), newRsp.GetGender())
}

func TestUserServer_CheckPassword(t *testing.T) {
	rawPassword := test_util.RandomString(10)
	slat, encryptedPassword := password.Encode(rawPassword, test_util.Options)
	encryptedPassword = fmt.Sprintf("$pbkdf2-sha512$%s$%s", slat, encryptedPassword)

	request := proto.PasswordCheckInfo{
		Password:          rawPassword,
		EncryptedPassword: encryptedPassword,
	}
	checkPasswordResponse, err := userClient.CheckPassword(context.Background(), &request)
	require.NoError(t, err)
	require.NotNil(t, checkPasswordResponse)
	require.True(t, checkPasswordResponse.Success)
}
