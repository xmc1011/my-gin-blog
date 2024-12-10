package repo

import (
	"errors"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"my-blog/internal/database/model"
	"my-blog/internal/schema"
	"time"
)

// 转换为Schema用户对象
func UserAuthModelToSchema(a *model.UserAuth) *schema.UserAuth {
	dest := schema.UserAuth{}
	copier.Copy(&dest, a)
	return &dest
}

// 用户权限存储
type RepoUserAuth struct {
	db *gorm.DB
}

// Get gorm.DB.Model
func (a *RepoUserAuth) GetDB() *gorm.DB {
	return a.db.Model(new(model.UserAuth))
}

// Get 查询指定数据
func (a *RepoUserAuth) GetUserAuthInfoByName(name string) (*schema.UserAuth, error) {
	userauth := new(model.UserAuth)

	result := a.GetDB().Where("username LIKE ?", name).First(&userauth)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return UserAuthModelToSchema(userauth), result.Error
}

func (a *RepoUserAuth) GetUserList(page, size int, loginType int8, nickname, username string) (list schema.UserAuthList, total int64, err error) {
	if loginType != 0 {
		a.db = a.GetDB().Where("login_type=?", loginType)
	}
	if username != "" {
		a.db = a.GetDB().Where("username LIKE?", "%"+username+"%")
	}
	//查询数据库
	// 执行查询并处理分页和总数
	result := a.db.Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
		Where("user_info.nickname LIKE ?", "%"+nickname+"%").
		Preload("UserInfo"). // 确保与模型字段一致
		Preload("Roles").
		Count(&total).                // 获取总记录数
		Scopes(Paginate(page, size)). // 应用分页
		Find(&list)                   // 执行查询并填充结果

	return list, total, result.Error
}

func (a *RepoUserAuth) UpdateUserAuthLoginInfo(id int, ipAddress, ipSource string) error {
	now := time.Now()
	userAuth := model.UserAuth{
		IpAddress:     ipAddress,
		IpSource:      ipSource,
		LastLoginTime: &now,
	}
	result := a.GetDB().Where("id", id).Updates(userAuth)
	return result.Error
}

func (a *RepoUserAuth) GetUserAuthInfoById(id int) (*schema.UserAuth, error) {
	// 查询 UserAuth 的基础数据
	userAuth := new(model.UserAuth)
	result := a.GetDB().Where("id = ?", id).First(userAuth)
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询关联的 Roles 数据
	var roles schema.RoleList
	roleResult := a.db.Table("user_auth_role").
		Select("role.*").
		Joins("JOIN role ON role.id = user_auth_role.role_id").
		Where("user_auth_role.user_auth_id = ?", id).
		Find(&roles)

	if roleResult.Error != nil {
		return nil, roleResult.Error
	}

	// 查询关联的 UserInfo 数据
	var userInfo model.UserInfo
	userInfoResult := a.db.Where("id = ?", userAuth.UserInfoId).First(&userInfo)
	if userInfoResult.Error != nil {
		return nil, userInfoResult.Error
	}

	// 转换为 schema.UserAuth
	schemaUserAuth := &schema.UserAuth{
		Model:         userAuth.Model,
		Username:      userAuth.Username,
		Password:      userAuth.Password,
		LoginType:     userAuth.LoginType,
		IpAddress:     userAuth.IpAddress,
		IpSource:      userAuth.IpSource,
		LastLoginTime: userAuth.LastLoginTime,
		IsDisable:     userAuth.IsDisable,
		IsSuper:       userAuth.IsSuper,
		UserInfoId:    userAuth.UserInfoId,
		UserInfo:      &userInfo,
		Roles:         roles,
	}

	return schemaUserAuth, nil
}
