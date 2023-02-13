## 工程说明

**OPS**运维平台公共包

- **jwt**
  -  生成,解析jwt Token 方法
```go
// 示例代码
// 可传入 Opthion 
// jwt.OptionWithExpireTime()
// jwt.OptionWithJwtSecret()
t := NewToken()
token,err := t.GenerateToken(userid)// 传入userid 返回token,err
claims,err := t.ParseToken(token) // 传入token解析  返回反序列化,及err

-------------------
// Gin 使用
r := g.New()
r.Use(jwt.JWTMiddleware()) // 注册中间件
r.Run()

  ```

- **Respoes**
  - 返回统一结构体
```go
// 示例
r := g.New()
r.GET("/ping", func(c *gin.Context) {
    response.Resp(c,response.SuccessCode,"",response.&Data{"Total":1,"Items":["1"]})
})

```

- **log**
- 日志统一输出
  -  输出路径统一为`pwd`/logs/
```go
g := gin.New()
g.Use(log.TraceIDMiddleware()) // 注册中间件

g.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong")
    log.DeaultLogs.Log.Info("aaa")
})
```

- **Client**

| 调用第三方服务(除了自身以外), 集成log trace_id 跨服务调用串联

eg:

```go
package main

import (
	`net/http`

	`github.com/gin-gonic/gin`

	`github.com/Jimi-Public/ops-common/client`
	`github.com/Jimi-Public/ops-common/log`
)

func main() {
	g := gin.New()
	g.Use(log.TraceIDMiddleware())

	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
		log.DeaultLogs.Log.Info("aaa")
		client.DefaultRequestBuilder.
			SetMethod(http.MethodGet).
			SetUrl("http://127.0.0.1:8081/ping").
			SetContext(c).Build(c) // 调用其他服务(接口)
	})
	go gin2()
	g.Run()
}

func gin2() {
	g := gin.New()
	g.Use(log.TraceIDMiddleware())

	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong gin2 ")
		log.DeaultLogs.Log.Info("aaa gin2")
	})

	g.Run(":8081")
}

/* 日志输出 trace_id 一致
{"file":"D:/workespaces/jimi/ops-common/main.go:26","func":"main.main.func1","level":"info","message":"aaa","time":"2023-02-13 19:59:34.397","trace_id":"64569c6a-6e2f-4f30-bd69-fde105b49c92"}

{"file":"D:/workespaces/jimi/ops-common/main.go:39","func":"main.gin2.func1","level":"info","message":"aaa gin2","time":"2023-02-13 19:59:36.455","trace_id":"64569c6a-6e2f-4f30-bd69-fde105b49c92"}


*/
```