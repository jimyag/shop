package paseto

import (
	"errors"
	"time"
)

// 验证token时返回的错误
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

//
// Payload
//  @Description: token认证的载荷
//
type Payload struct {
	IssuedAt  time.Time
	ExpiredAt time.Time
	UID       int32
	Role      int32
}

//
// NewPayload
//  @Description:
//  @param uid
//  @param role
//  @return *Payload
//  @return error
//
func NewPayload(uid int32, role int32) (*Payload, error) {
	payload := &Payload{
		UID:  uid,
		Role: role,
	}
	return payload, nil
}

//
// Valid
//  @Description: 检查token是否过期
//  @receiver payload
//  @return error
//
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
