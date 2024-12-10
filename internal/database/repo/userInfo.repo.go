package repo

import (
	"gorm.io/gorm"
	"my-blog/internal/database/model"
)

// 用户权限存储
type RepoUserInfo struct {
	db *gorm.DB
}

// Get gorm.DB.Model
func (a *RepoUserInfo) GetDB() *gorm.DB {
	return a.db.Model(new(model.UserInfo))
}

func (a *RepoUserInfo) GetUserInfoById(id int) (*model.UserInfo, error) {
	var userInfo model.UserInfo
	result := a.GetDB().Where("id", id).First(&userInfo)
	return &userInfo, result.Error
}

func (a *RepoUserInfo) UpdateUserInfo(id int, nickname, avater, intro, website string) error {
	userInfo := model.UserInfo{
		Model:    model.Model{ID: id},
		Nickname: nickname,
		Avatar:   avater,
		Intro:    intro,
		Website:  website,
	}
	result := a.GetDB().Select("nickname", "avater", "intro", "website").Updates(userInfo)
	return result.Error
}
