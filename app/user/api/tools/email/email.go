package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"

	"github.com/jimyag/shop/app/user/api/global"
)

//
// Send
//  @Description: 发送邮件
//  @param to 发送给谁，接收方的邮箱的
//  @param subject 标题
//  @param body 内容
//  @return error
//
func Send(to []string, subject string, body string) error {
	from := global.RemoteConfig.Email.Form
	nickname := global.RemoteConfig.Email.Nickname
	secret := global.RemoteConfig.Email.Secret
	host := global.RemoteConfig.Email.Host
	port := global.RemoteConfig.Email.Port
	isSSL := global.RemoteConfig.Email.IsSsl

	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
