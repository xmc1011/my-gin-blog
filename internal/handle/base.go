package handle

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-blog/internal/global"
	"my-blog/internal/schema"
	"my-blog/internal/service"
	"my-blog/internal/utils/logger"
	"net/http"
)

/*
响应设计方案：不使用 HTTP 码来表示业务状态, 采用业务状态码的方式
- 只要能到达后端的请求, HTTP 状态码都为 200
- 业务状态码为 0 表示成功, 其他都表示失败
- 当后端发生 panic 并且被 gin 中间件捕获时, 才会返回 HTTP 500 状态码
*/

// 响应结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 业务状态码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // 响应数据
}

// HTTP 码 + 业务码 + 消息 + 数据
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// 业务码 + 数据
func ReturnResponse(c *gin.Context, r global.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

// 成功业务码 + 数据
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, global.OkResult, data)
}

// 所有可预料的错误 = 业务错误 + 系统错误, 在业务层面处理, 返回 HTTP 200 状态码
// 对于不可预料的错误, 会触发 panic, 由 gin 中间件捕获, 并返回 HTTP 500 状态码
// err 是业务错误, data 是错误数据 (可以是 error 或 string)
func ReturnError(c *gin.Context, r global.Result, data any) {
	logrus.Info("[Func-ReturnError] " + r.Msg())

	var val string = r.Msg()
	if data != nil {
		switch v := data.(type) {
		case error:
			val = v.Error()
		case string:
			val = v
		}
		logrus.Error(val)
	}
	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Msg(),
			Data:    val,
		},
	)
}

// 获取 *gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(global.CTX_DB).(*gorm.DB)
}

// 获取 *redis.Client
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(global.CTX_RDB).(*redis.Client)
}

/*
获取当前登录用户信息
1. 能从 gin Context 上获取到 user 对象, 说明本次请求链路中获取过了
2. 从 session 中获取到 uid
3. 根据 uid 获取用户信息, 并设置到 gin Context 上
*/
func CurrentUserAuth(c *gin.Context) (*schema.UserAuth, error) {
	key := global.CTX_USER_AUTH

	//1
	if cache, exist := c.Get(key); exist && cache != nil {
		logger.Debugf("[Func-CurrentUserAuth] get from cache: " + cache.(*schema.UserAuth).Username)
		return cache.(*schema.UserAuth), nil
	}

	//2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	//3
	user, err := service.GetUserAuthInfoById(id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
