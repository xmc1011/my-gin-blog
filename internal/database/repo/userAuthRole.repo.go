package repo

import (
	"gorm.io/gorm"
	"my-blog/internal/database/model"
)

// RepoUserAuthRole 用户权限角色存储
type RepoUserAuthRole struct {
	db *gorm.DB
}

// GetDB 返回 gorm.DB 对象，用于模型操作
func (a *RepoUserAuthRole) GetDB() *gorm.DB {
	return a.db.Model(new(model.UserAuthRole))
}

// GetRoleIdsByUserId 通过 user_auth_id 获取 role_id 列表
func (a *RepoUserAuthRole) GetRoleIdsByUserId(userAuthID int) (ids []int, err error) {
	result := a.GetDB().Where("user_auth_id = ?", userAuthID).Pluck("role_id", &ids)
	return ids, result.Error
}
