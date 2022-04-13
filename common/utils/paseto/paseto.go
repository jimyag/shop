package paseto

import (
	"crypto/ed25519"
	"time"

	"github.com/o1egl/paseto"
)

//
// PasetoMaker
//  @Description: is a PASETO token maker
//
type PasetoMaker struct {
	pastor     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	duration   time.Duration
}

//
// NewPasetoMaker
//  @Description: creates a new PasetoMaker
//  @param privateKey 私钥
//  @param publicKey 公钥
//  @param duration 有效时间 小时
//  @return *PasetoMaker
//  @return error
//
func NewPasetoMaker(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey, duration time.Duration) (*PasetoMaker, error) {
	maker := &PasetoMaker{
		pastor:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
		duration:   duration,
	}
	return maker, nil
}

//
// CreateToken
//  @Description: 使用载荷生成token
//  @receiver maker
//  @param payload
//  @return string
//  @return error
//
func (maker *PasetoMaker) CreateToken(payload *Payload) (string, error) {
	payload.IssuedAt = time.Now()
	payload.ExpiredAt = time.Now().Add(time.Hour * maker.duration)
	token, err := maker.pastor.Sign(maker.privateKey, payload, nil)
	return token, err
}

//
// VerifyToken
//  @Description: 验证token是否合法
//  @receiver maker
//  @param token
//  @return *Payload
//  @return error
//
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.pastor.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if payload.Valid() != nil {
		return nil, err
	}
	return payload, nil
}
