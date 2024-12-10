package repo

import "gorm.io/gorm"

var (
	UserAuthRepo *RepoUserAuth
	UserInfoRepo *RepoUserInfo
	UserAuthRole *RepoUserAuthRole
)

func InitRepo(db *gorm.DB) {
	UserAuthRepo = &RepoUserAuth{db: db}
	UserInfoRepo = &RepoUserInfo{db: db}
	UserAuthRole = &RepoUserAuthRole{db: db}
}
