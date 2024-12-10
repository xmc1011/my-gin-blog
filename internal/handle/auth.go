package handle

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"my-blog/internal/config"
	"my-blog/internal/database/model"
	"my-blog/internal/global"
	"my-blog/internal/service"
	"my-blog/internal/utils"
	"my-blog/internal/utils/jwt"
	"my-blog/internal/utils/logger"
	"strconv"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginVO struct {
	model.UserInfo

	// 点赞 Set: 用于记录用户点赞过的文章, 评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

// @Summary 登录
// @Description 登录
// @Tags UserAuth
// @Param form body LoginReq true "登录"
// @Accept json
// @Produce json
// @Success 200
// @Router /login [post]
func (*UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	rdb := GetRDB(c)

	userAuth, err := service.GetUserAuthInfoByName(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, global.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	// 检查密码是否正确
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, global.ErrPassword, nil)
		return
	}

	// 获取IP相关信息
	ipAddress := utils.IP.GetIpAddress(c)
	logger.Infof("ipAddress", ipAddress)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)
	logger.Infof("ipSource", ipSource)

	// 获取用户信息
	userInfo, err := service.GetUserInfoById(userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, global.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	// 获取角色ID
	roleIds, err := service.GetRoleIdsByUserId(userAuth.ID)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	articleLikeSet, err := rdb.SMembers(rctx, global.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(rctx, global.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	//生成TOKEN
	conf := config.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		ReturnError(c, global.ErrTokenCreate, err)
		return
	}
	logger.Infof("token", token)

	// 更新用户验证信息: ip 信息 + 上次登录时间
	err = service.UpdateUserAuthLoginInfo(userAuth.ID, ipAddress, ipSource)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	logger.Infof("用户登录成功: " + userAuth.Username)

	session := sessions.Default(c)
	session.Set(global.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	//删除 Redis 中的离线状态
	offlineKey := global.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	ReturnSuccess(c, LoginVO{
		UserInfo:       *userInfo,
		Token:          token,
		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
	})
}

// @Summary 退出登录
// @Description 退出登录
// @Tags UserAuth
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /logout [post]
func (*UserAuth) Logout(c *gin.Context) {
	c.Set(global.CTX_USER_AUTH, nil)
	// 已经退出登录

	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}
	session := sessions.Default(c)
	session.Delete(global.CTX_USER_AUTH)
	session.Save()

	// 删除 Redis 中的在线状态
	rdb := GetRDB(c)
	onlineKey := global.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rctx, onlineKey)
	ReturnSuccess(c, nil)

}
