package validate

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

//
// Validate
//  @Description: 验证传入的参数是否符合要求
//  			  传入的只能是 struct
//  @param data
//  @param validate
//  @param trans
//  @return interface{}
//  @return error
//
func Validate(data interface{}, validate *validator.Validate, trans ut.Translator) (interface{}, error) {
	err := validate.Struct(data)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			errMsg := make([]interface{}, 0)
			for _, fieldError := range errs {
				errMsg = append(errMsg, fieldError.Translate(trans))
			}
			return errMsg, err

		}
	}
	return nil, nil
}
