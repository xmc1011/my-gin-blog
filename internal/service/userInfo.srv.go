package service

import (
	"my-blog/internal/database/model"
	"my-blog/internal/database/repo"
)

// GetUserInfoById 查询指定用户名的用户信息
func GetUserInfoById(id int) (*model.UserInfo, error) {
	// 调用 Repo 层的 GetUserAuthInfoByName 方法
	return repo.UserInfoRepo.GetUserInfoById(id)
}
