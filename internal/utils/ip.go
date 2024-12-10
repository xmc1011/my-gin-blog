package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/mssola/user_agent" // 替换为 mssola/user_agent
	"my-blog/internal/utils/logger"
	"net"
	"strings"
)

var IP = new(ipUtil)

type ipUtil struct{}

// 获取用户发送请求的 IP 地址
func (*ipUtil) GetIpAddress(c *gin.Context) (ipAddress string) {
	// c.ClientIP() 获取的是代理服务器的 IP (Nginx)

	// X-Real-IP: Nginx 服务代理, 本项目明确使用 Nginx 作代理, 因此优先获取这个
	ipAddress = c.Request.Header.Get("X-Real-IP")

	// X-Forwarded-For 经过 HTTP 代理或 负载均衡服务器时会添加该项
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ips := c.Request.Header.Get("X-Forwarded-For")
		splitIps := strings.Split(ips, ",")
		if len(splitIps) > 0 {
			ipAddress = splitIps[0]
		}
	}

	// Pdoxy-Client-IP: Apache 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("Proxy-Client-IP")
	}

	// WL-Proxy-Client-IP: Weblogic 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
	}

	// RemoteAddr: 发出请求的远程主机的 IP 地址 (经过代理会设置为代理机器的 IP)
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.RemoteAddr
	}

	// 检测到是本机 IP, 读取其局域网 IP 地址
	if strings.HasPrefix(ipAddress, "127.0.0.1") || strings.HasPrefix(ipAddress, "[::1]") {
		ip, err := externalIP()
		if err != nil {
			logger.Errorf("GetIpAddress, externalIP, err: ", err)
		}
		ipAddress = ip.String()
	}

	// 如果有多个 IP，取第一个
	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Index(ipAddress, ",") > 0 {
			ipAddress = ipAddress[:strings.Index(ipAddress, ",")]
		}
	}
	return ipAddress
}

// 获取 IP 来源
var vIndex []byte // 缓存 VectorIndex 索引, 减少一次固定的 IO 操作

func (*ipUtil) GetIpSource(ipAddress string) string {
	var dbPath = "/home/xmc/study/my-blog/assets/ip2region.xdb"
	if vIndex == nil {
		var err error
		vIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
		if err != nil {
			logger.Errorf(fmt.Sprintf("failed to load vector index from `%s`: %s\n", dbPath, err))
			return ""
		}
	}
	searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex)

	if err != nil {
		logger.Errorf("failed to create searcher with vector index: ", err)
		return ""
	}
	defer searcher.Close()

	// 查询 IP 地理信息
	region, err := searcher.SearchByStr(ipAddress)
	if err != nil {
		logger.Errorf(fmt.Sprintf("failed to search ip(%s): %s\n", ipAddress, err))
		return ""
	}
	return region
}

// 获取 IP 简易信息
func (i *ipUtil) GetIpSourceSimpleIdle(ipAddress string) string {
	region := i.GetIpSource(ipAddress)

	// 检测到是内网, 直接返回 "内网IP"
	if strings.Contains(region, "内网IP") {
		return "内网IP"
	}

	// 中国|0|江苏省|苏州市|电信
	ipSource := strings.Split(region, "|")
	if ipSource[0] != "中国" && ipSource[0] != "0" {
		return ipSource[0]
	}
	if ipSource[2] == "0" {
		ipSource[2] = ""
	}
	if ipSource[3] == "0" {
		ipSource[3] = ""
	}
	if ipSource[4] == "0" {
		ipSource[4] = ""
	}
	if ipSource[2] == "" && ipSource[3] == "" && ipSource[4] == "" {
		return ipSource[0]
	}
	return ipSource[2] + ipSource[3] + " " + ipSource[4]
}

// 获取用户代理信息
func (*ipUtil) GetUserAgent(c *gin.Context) *user_agent.UserAgent {
	ua := c.Request.UserAgent()
	return user_agent.New(ua)
}

// 获取非 127.0.0.1 的局域网 IP
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
