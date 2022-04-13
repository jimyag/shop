package initialize

import (
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/goods/api/global"
)

//
// InitValidateAndTrans
//  @Description: 初始化 validate 和 trans
//
func InitValidateAndTrans() {
	global.Validate = validator.New()
	uni := ut.New(zh.New(), zh.New())
	var ok bool
	global.Trans, ok = uni.GetTranslator("zh")
	if !ok {
		global.Logger.Error("获得翻译器失败")
	}

	err := zhtranslations.RegisterDefaultTranslations(global.Validate, global.Trans)
	if err != nil {
		global.Logger.Error("注册翻译器失败", zap.Error(err))
	}

	global.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	global.Logger.Info("翻译器注册成功......")
}
