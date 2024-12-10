package test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	c "my-blog/internal/config"
	"my-blog/internal/utils/logger"
	"testing"
)

func TestInitLogger(t *testing.T) {

	conf := c.ReadConfig("../../../config") //读取配置文件

	// 调用 InitLogger 初始化 logrus 实例
	logger := logger.InitLogger(conf)

	// 检查日志级别是否正确设置
	assert.Equal(t, logrus.DebugLevel, logger.GetLevel(), "日志级别不正确")

	logger.Debug("this is debug log")
	logger.Info("this is info log")
}
