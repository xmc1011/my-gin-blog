package service

import (
	"my-blog/internal/database/repo"
)

// GetRoleIdsByUserId 查询指定用户名的用户信息
func GetRoleIdsByUserId(userAuthID int) ([]int, error) {
	// 调用 Repo 层的 GetUserAuthInfoByName 方法
	return repo.UserAuthRole.GetRoleIdsByUserId(userAuthID)
}
