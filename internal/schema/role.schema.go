package schema

import "my-blog/internal/database/model"

type Role struct {
	model.Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources ResourceList `json:"resources" gorm:"many2many:role_resource"`
	Menus     MenuList     `json:"menus" gorm:"many2many:role_menu"`
	Users     UserAuthList `json:"users" gorm:"many2many:user_auth_role"`
}

type RoleList []*Role //涉及返回多个结构体进行封装
