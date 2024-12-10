package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"my-blog/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 测试 GetIpAddress 方法
func TestGetIpAddress(t *testing.T) {
	// 创建一个模拟的 gin.Context
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.Header("X-Real-IP", "192.168.118.131")
		ipAddress := utils.IP.GetIpAddress(c)
		assert.Equal(t, "192.168.118.131", ipAddress)
	})

	// 模拟请求
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

// 测试 GetIpSource 方法
func TestGetIpSource(t *testing.T) {
	ipAddress := "8.8.8.8" // 你可以使用任意有效的 IP 地址来模拟查询
	ipSource := utils.IP.GetIpSource(ipAddress)
	assert.NotEmpty(t, ipSource)
	// 可以根据具体的 IP 数据库进行更详细的测试，或者用 mock 数据
}

// 测试 GetIpSourceSimpleIdle 方法
func TestGetIpSourceSimpleIdle(t *testing.T) {
	ipAddress := "8.8.8.8" // 使用任意 IP 地址进行测试
	ipSourceSimple := utils.IP.GetIpSourceSimpleIdle(ipAddress)
	assert.NotEmpty(t, ipSourceSimple)
	// 可以根据需要自定义匹配规则
}

// 测试 GetUserAgent 方法
func TestGetUserAgent(t *testing.T) {
	// 创建一个模拟的 gin.Context
	r := gin.Default()
	r.GET("/test-user-agent", func(c *gin.Context) {
		// 模拟请求头
		c.Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		ua := utils.IP.GetUserAgent(c)
		assert.NotNil(t, ua)
		assert.Equal(t, "Windows NT 10.0; Win64; x64", ua.OS())
	})

	// 模拟请求
	req, err := http.NewRequest("GET", "/test-user-agent", nil)
	assert.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

// 测试 GetIpAddress 从 X-Forwarded-For 中获取
func TestGetIpAddressFromXForwardedFor(t *testing.T) {
	// 创建一个模拟的 gin.Context
	r := gin.Default()
	r.GET("/test-forwarded", func(c *gin.Context) {
		c.Header("X-Forwarded-For", "192.168.1.2, 192.168.1.1")
		ipAddress := utils.IP.GetIpAddress(c)
		assert.Equal(t, "192.168.1.2", ipAddress)
	})

	// 模拟请求
	req, err := http.NewRequest("GET", "/test-forwarded", nil)
	assert.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}
