package test

import (
	"github.com/stretchr/testify/assert"
	c "my-blog/internal/config"
	"my-blog/internal/utils"
	"testing"
)

// 测试数据库连接
func TestInitRedis(t *testing.T) {
	conf := c.ReadConfig("../../../config") //读取配置文件
	// 使用测试配置初始化数据库
	rdb := utils.InitRedis(conf)

	// 验证数据库连接是否成功
	assert.NotNil(t, rdb, "redis连接失败")

}
