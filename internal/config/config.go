package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// 存储服务器配置
type Server struct {
	Mode          string // debug | release
	Port          string
	DbType        string // mysql | sqlite
	DbAutoMigrate bool   // 是否自动迁移数据库表结构
	DbLogMode     string // silent | error | warn | info
}

// 存储JWT配置
type JWT struct {
	Secret string
	Expire int64 // hour
	Issuer string
}

// 存储 MySQL 配置
type Mysql struct {
	Host     string // 服务器地址
	Port     string // 端口
	Config   string // 高级配置
	Dbname   string // 数据库名
	Username string // 数据库用户名
	Password string // 数据库密码
}

// 存储 SQLite 配置
type Sqlite struct {
	Dsn string // Data Source Name
}

// 存储 Redis 配置
type Redis struct {
	DB       int    // 指定 Redis 数据库
	Addr     string // 服务器地址:端口
	Password string // 密码
}

// 存储会话配置
type Session struct {
	Name   string
	Salt   string
	MaxAge int
}

// 存储日志配置
type Log struct {
	Level     string // debug | info | warn | error
	Prefix    string
	Format    string // text | json
	Directory string
}

// 存储邮件配置
type Email struct {
	From     string // 发件人 要发邮件的邮箱
	Host     string // 服务器地址, 例如 smtp.qq.com 前往要发邮件的邮箱查看其 smtp 协议
	Port     int    // 前往要发邮件的邮箱查看其 smtp 协议端口, 大多为 465
	SmtpPass string // 邮箱密钥 不是密码是开启smtp后给你的密钥
	SmtpUser string // 邮箱账号
}

// 存储验证码配置
type Captcha struct {
	SendEmail  bool // 是否通过邮箱发送验证码
	ExpireTime int  // 过期时间
}

// 存储文件上传配置
type Upload struct {
	// Size      int    // 文件上传的最大值
	OssType   string // local | qiniu
	Path      string // 本地文件访问路径
	StorePath string // 本地文件存储路径
}

// 存储七牛云配置
type Qiniu struct {
	ImgPath       string // 外链链接
	Zone          string // 存储区域
	Bucket        string // 空间名称
	AccessKey     string // 秘钥AK
	SecretKey     string // 秘钥SK
	UseHTTPS      bool   // 是否使用https
	UseCdnDomains bool   // 上传是否使用 CDN 上传加速
}

// Config 用于将所有配置项组合在一起
type Config struct {
	Server  Server
	JWT     JWT
	Mysql   Mysql
	Sqlite  Sqlite
	Redis   Redis
	Session Session
	Log     Log
	Email   Email
	Captcha Captcha
	Upload  Upload
	Qiniu   Qiniu
}

var Conf *Config

// 获取配置
func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

// 使用viper获取配置文件
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigName("config") // 配置文件名，不包含扩展名
	v.SetConfigType("yaml")   // 配置文件类型
	v.AddConfigPath(path)     // 配置文件路径
	v.AutomaticEnv()          //允许环境变量
	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败：" + err.Error())
	}
	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}

	log.Println("配置文件内容加载成功: ", path)
	return Conf
}

// 数据库类型
func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

// 数据库链接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		conf := Conf.Mysql
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config)
	case "sqlite":
		return Conf.Sqlite.Dsn
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.Sqlite.Dsn == "" {
			Conf.Sqlite.Dsn = "sqlite"
		}
		return Conf.Sqlite.Dsn
	}

}
