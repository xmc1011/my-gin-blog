package global

import "fmt"

const (
	SUCCESS = 0   // 成功业务码
	FAIL    = 500 // 失败业务码
)

// 自定义业务 error 类型
type Result struct {
	code int
	msg  string
}

func (e Result) Code() int {
	return e.code
}

func (e Result) Msg() string {
	return e.msg
}

var (
	OkResult   = RegisterResult(SUCCESS, "OK")
	FailResult = RegisterResult(FAIL, "FAIL")
)

var (
	_codes    = map[int]struct{}{}   // 注册过的错误码集合, 防止重复
	_messages = make(map[int]string) // 根据错误码获取错误信息
)

// 注册一个响应码, 不允许重复注册
func RegisterResult(code int, msg string) Result {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	if msg == "" {
		panic("错误码消息不能为空")
	}

	_codes[code] = struct{}{}
	_messages[code] = msg

	return Result{
		code: code,
		msg:  msg,
	}
}

var (
	ErrRequest = RegisterResult(9001, "请求参数格式错误")
	ErrDbOp    = RegisterResult(9004, "数据库操作异常")
	ErrRedisOp = RegisterResult(9005, "Redis 操作异常")

	ErrPassword     = RegisterResult(1002, "密码错误")
	ErrUserNotExist = RegisterResult(1003, "该用户不存在")

	ErrTokenNotExist    = RegisterResult(1201, "TOKEN 不存在，请重新登陆")
	ErrTokenRuntime     = RegisterResult(1202, "TOKEN 已过期，请重新登陆")
	ErrTokenWrong       = RegisterResult(1203, "TOKEN 不正确，请重新登陆")
	ErrTokenType        = RegisterResult(1204, "TOKEN 格式错误，请重新登陆")
	ErrTokenCreate      = RegisterResult(1205, "TOKEN 生成失败")
	ErrPermission       = RegisterResult(1206, "权限不足")
	ErrForceOffline     = RegisterResult(1207, "您已被强制下线")
	ErrForceOfflineSelf = RegisterResult(1208, "不能强制下线自己")
)
