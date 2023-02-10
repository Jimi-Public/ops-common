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
claims,err := t.ParseToken(token) # 传入token解析  返回反序列化,及err

  ```

- ***Respoes**
  - 返回统一结构体
```go
// 示例
r := g.New()
r.GET("/ping", func(c *gin.Context) {
    response.Resp(c,response.SuccessCode,"",response.&Ddata{"Total":1,"Items":["1"]})
})

```

- **log**
- 日志统一输出
  -  输出路径统一为`pwd`/logs/
```go
g := gin.New()
g.Use(log.TraceIDMiddleware())

g.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong")
    log.DeaultLogs.Log.Info("aaa")
})
```