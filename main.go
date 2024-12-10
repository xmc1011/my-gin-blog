package main

import (
	"github.com/gin-gonic/gin"
	"my-blog/internal"
	c "my-blog/internal/config"
	"my-blog/internal/database"
	"my-blog/internal/middleware"
	"my-blog/internal/utils"
	"my-blog/internal/utils/logger"
)

func main() {

	conf := c.ReadConfig("./config") //读取配置文件

	_ = logger.InitLogger(conf)
	database.InitGormDB(conf)
	rdb := utils.InitRedis(conf)

	//初始化gin
	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})
	// 开发模式使用 gin 自带的日志和恢复中间件, 生产模式使用自定义的中间件
	if conf.Server.Mode == "debug" {
		r.Use(gin.Logger(), gin.Recovery()) // gin 自带的日志和恢复中间件, 挺好用的
	} else {
		r.Use(middleware.Recovery(true), middleware.Logger())
	}
	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithRedisDB(rdb))
	internal.RegisterHandlers(r)
	serverPort := conf.Server.Port
	r.Run(serverPort)
}
