package auth_code

import (
	"math/rand"
	"strconv"
	"time"
)

//
// Code
//  @Description:生成5位数字验证码
//  @return int
//
func Code() int {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(89999) + 10000
	return code
}

//
// CodeBody
//  @Description:发送给用户邮件的模板
//  @param code
//  @return string
//
func CodeBody(code int, content string) string {
	header := `<h1>【管理系统】 您正在` + content + ` 。您的验证码是：`
	return header + strconv.Itoa(code) + `，请勿向他人泄露，转发可能导致账号被盗。如非本人操作，可忽略本消息。5分钟内有效</h1>`
}
