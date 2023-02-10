/*
@Project: user-manager
@Author:  WangChaoQun
@Date:    2023/2/7
@IDE:     GoLand
@File:    response.go
*/

package response

import "github.com/gin-gonic/gin"

type Code int

type Response struct {
	Code    Code   `json:"code" `    // 业务状态码
	Message string `json:"message" ` // 提示信息
	Data    Data   `json:"data" `    // 任何数据
}

func newResponse(code Code, message string, data Data) *Response {
	return &Response{Code: code, Message: message, Data: data}
}

type Data struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

const (
	SuccessCode   Code = 20000 // 成功状态码
	ForbiddenCode Code = 21000 // 操作拒绝
	FailCode      Code = 22000 // 操作失败 后端异常
	AuthFail      Code = 40300 // 认证失败
)

func Resp(ctx *gin.Context, code Code, message string, data Data) {
	r := newResponse(code, message, data)
	ctx.JSON(200, r)
}
