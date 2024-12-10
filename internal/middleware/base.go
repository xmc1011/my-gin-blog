package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"my-blog/internal/global"
	"net/http"
	"time"
)

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_RDB).(*redis.Client) 来使用
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_RDB, rdb)
		ctx.Next()
	}
}

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_DB).(*gorm.DB) 来使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_DB, db)
		ctx.Next()
	}
}

// WithCookieStore 基于 cookie 的 session
func WithCookieStore(name, secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}

// CORS 跨域请求
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}

// Logger 日志记录
func Logger() gin.HandlerFunc {
	logger := logrus.New()
	return func(c *gin.Context) {
		start := time.Now()       //请求开始时间
		c.Next()                  //继续执行请求
		cost := time.Since(start) //请求耗时
		//通过 WithFields 为日志添加更多的上下文信息
		logger.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"status": c.Writer.Status(),
			"method": c.Request.Method,
			"ip":     c.ClientIP(),
			"size":   c.Writer.Size(),
			"cost":   cost,
			// "body":   c.Request.PostForm.Encode(), // 如果需要记录请求体
			// "model-agent": c.Request.UserAgent(),
			// "errors":    c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}).Info("[GIN]") // 使用 Info 级别记录日志
	}
}

// 恢复中间件
// Recovery 恢复中间件
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 返回 500 内部服务器错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
