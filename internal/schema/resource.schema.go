package schema

import "my-blog/internal/database/model"

type Resource struct {
	model.Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles RoleList `json:"roles" gorm:"many2many:role_resource"`
}

type ResourceList []*Resource //涉及返回多个结构体进行封装
