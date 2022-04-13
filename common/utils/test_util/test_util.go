package test_util

import (
	"crypto/sha512"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
)

const (
	alphabet = "abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ0123456789"
)

var (
	Options = &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())
}

//
// RandomString
//  @Description: 随机指定长度的字符串 包含大小写字母和数字
//  @param n
//  @return string
//
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//
// RandomInt
//  @Description: 生成指定范围的int
//  @param min
//  @param max
//  @return int64
//
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//
// RandomEmail
//  @Description: 随机生成 "jimyag%s@126.com"格式的邮箱 len(s) = 3
//  @return string
//
func RandomEmail() string {
	return fmt.Sprintf("jimyag%s@126.com", RandomString(3))
}

type Password struct {
	RawPassword       string
	Slat              string
	EncryptedPassword string
}

//
// RandomPassword
//  @Description:
//  @return p
//
func RandomPassword() (p Password) {
	rawPassword := RandomString(10)
	slat, encryptedPassword := password.Encode(rawPassword, Options)
	p = Password{
		RawPassword:       rawPassword,
		Slat:              slat,
		EncryptedPassword: fmt.Sprintf("$pbkdf2-sha512$%s$%s", slat, encryptedPassword),
	}
	return
}

//
// RandomNickName
//  @Description: 随机昵称 jimyag%s len(s) = 5
//  @return string
//
func RandomNickName() string {
	return fmt.Sprintf("jimyag%s", RandomString(5))
}

//
// RandomGender
//  @Description: 随机性别 "male", "female", "middle"
//  @return string
//
func RandomGender() string {
	gender := []string{"male", "female", "middle"}
	n := len(gender)
	return gender[rand.Intn(n)]
}

//
// RandomFloat
//  @Description: 求指定范围的随机
//  @param min
//  @param max
//  @return float32
//
func RandomFloat(min, max float32) float32 {
	return min + rand.Float32()*max - (max - min)
}

func RandomPrice() float32 {
	return RandomFloat(10.0, 2000.0)
}
