package schema

import (
	"my-blog/internal/database/model"
	"time"
)

type UserAuth struct {
	model.Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置
	UserInfoId    int        `json:"user_info_id"`

	UserInfo *model.UserInfo `json:"info"`
	Roles    []*model.Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

type UserAuthList []*UserAuth //涉及返回多个结构体进行封装
