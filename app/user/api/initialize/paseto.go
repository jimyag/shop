package initialize

import (
	"crypto/ed25519"
	"encoding/hex"
	"time"

	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/global"
	"github.com/jimyag/shop/app/user/api/tools/paseto"
)

//
// InitPaseto
//  @Description: 初始化paseto
//
func InitPaseto() {
	b, err := hex.DecodeString(global.RemoteConfig.Secret.PrivateKey)
	if err != nil {
		global.Logger.Error("初始化秘钥错误", zap.Error(err))
	}
	privateKey := ed25519.PrivateKey(b)

	b, err = hex.DecodeString(global.RemoteConfig.Secret.PrivateKey)
	if err != nil {
		global.Logger.Error("初始化公钥错误", zap.Error(err))
	}
	publicKey := ed25519.PublicKey(b)

	global.PasetoMaker, _ = paseto.NewPasetoMaker(privateKey, publicKey, time.Duration(global.RemoteConfig.Secret.Duration))
	global.Logger.Info("初始化PASETO成功......")
}
