package service

import (
	"my-blog/internal/database/repo"
	"my-blog/internal/schema"
)

// GetUserAuthInfoByName 查询指定用户名的用户认证信息
func GetUserAuthInfoByName(name string) (*schema.UserAuth, error) {
	// 调用 Repo 层的 GetUserAuthInfoByName 方法
	return repo.UserAuthRepo.GetUserAuthInfoByName(name)
}

// GetUserAuthInfoByName 查询指定用户名的用户认证信息
func GetUserAuthInfoById(id int) (*schema.UserAuth, error) {
	// 调用 Repo 层的 GetUserAuthInfoByName 方法
	return repo.UserAuthRepo.GetUserAuthInfoById(id)
}

// UpdateUserAuthLoginInfo 更新用户登陆信息
func UpdateUserAuthLoginInfo(id int, ipAddress, ipSource string) error {
	return repo.UserAuthRepo.UpdateUserAuthLoginInfo(id, ipAddress, ipSource)
}
