package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	c "my-blog/internal/config"
)

func InitRedis(conf *c.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.118.131:63799",
		Password: "66554321", // 没有密码，默认值
		DB:       0,          // 默认DB 0
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis 连接失败: ", err)
	}

	log.Println("Redis 连接成功", conf.Redis.Addr, conf.Redis.DB, conf.Redis.Password)
	return rdb
}
