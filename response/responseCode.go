/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/10
@IDE:     GoLand
@File:    responseCode.go
*/

package response

// 业务状态码
const (
	SuccessCode   Code = 20000 // 成功状态码
	ForbiddenCode Code = 21000 // 操作拒绝
	ParamsErr     Code = 23000 // 参数异常
	FailCode      Code = 22000 // 操作失败 后端异常
	AuthFail      Code = 40300 // 认证失败
	TokenExpire   Code = 40301 // token 过期 
)
