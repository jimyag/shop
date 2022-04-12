package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// Response
//  @Description: 统一的响应
//
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}

const (
	SUCCESS = 0
	ERROR   = 500
)

var codeMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "失败",
}

func getErrMsg(code int) string {
	return codeMsg[code]
}
func result(code int, data interface{}, msg interface{}, context *gin.Context) {
	context.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

//
// Ok
//  @Description: 成功
//  @param context
//
func Ok(context *gin.Context) {
	result(SUCCESS, nil, getErrMsg(SUCCESS), context)
}

//
// Fail
//  @Description:  失败
//  @param context
//
func Fail(context *gin.Context) {
	result(ERROR, nil, getErrMsg(ERROR), context)
}

//
// OkWithData
//  @Description: 成功并携带数据
//  @param data
//  @param context
//
func OkWithData(data interface{}, context *gin.Context) {
	result(SUCCESS, data, getErrMsg(SUCCESS), context)
}

//
// OkWithMsg
//  @Description: 成功携带msg
//  @param msg
//  @param context
//
func OkWithMsg(msg interface{}, context *gin.Context) {
	result(SUCCESS, nil, msg, context)
}

//
// OkWithDataMsg
//  @Description: 成功携带数据和msg
//  @param data
//  @param msg
//  @param context
//
func OkWithDataMsg(data interface{}, msg interface{}, context *gin.Context) {
	result(SUCCESS, data, msg, context)
}

//
// FailWithData
//  @Description: 失败携带数据
//  @param data
//  @param context
//
func FailWithData(data interface{}, context *gin.Context) {
	result(ERROR, data, getErrMsg(ERROR), context)
}

//
// FailWithMsg
//  @Description: 失败msg
//  @param msg
//  @param context
//
func FailWithMsg(msg interface{}, context *gin.Context) {
	result(ERROR, nil, msg, context)
}

//
// FailWithDataMsg
//  @Description: 失败数据和msg
//  @param data
//  @param msg
//  @param context
//
func FailWithDataMsg(data interface{}, msg interface{}, context *gin.Context) {
	result(ERROR, data, msg, context)
}
