package test

import (
	"github.com/stretchr/testify/assert"
	c "my-blog/internal/config"
	"my-blog/internal/database"
	"testing"
)

// 测试数据库连接
func TestInitDatabase(t *testing.T) {
	conf := c.ReadConfig("../../../config") //读取配置文件
	// 使用测试配置初始化数据库
	err := database.InitGormDB(conf)

	// 验证数据库连接是否成功
	assert.NoError(t, err, "Expected no error during database initialization")

}
