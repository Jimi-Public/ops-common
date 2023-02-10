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

func Resp(ctx *gin.Context, code Code, message string, data Data) {
	r := newResponse(code, message, data)
	ctx.JSON(200, r)
}
