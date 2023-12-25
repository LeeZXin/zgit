package apicommon

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"zgit/pkg/apicode"
	"zgit/pkg/apisession"
	"zgit/pkg/i18n"
	"zgit/standalone/modules/model/usermd"
)

const (
	LoginUser = "loginUser"
)

func CheckLogin(c *gin.Context) {
	sessionId := c.GetHeader("Authorization")
	if sessionId == "" {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.NotLoginCode.Int(),
			Message: i18n.GetByKey(i18n.SystemNotLogin),
		})
		c.Abort()
		return
	}
	sessionStore := apisession.GetStore()
	session, b, err := sessionStore.GetBySessionId(sessionId)
	if err != nil {
		logger.Logger.WithContext(c.Request.Context()).Error(err)
		c.JSON(http.StatusInternalServerError, ginutil.BaseResp{
			Code:    apicode.InternalErrorCode.Int(),
			Message: i18n.GetByKey(i18n.SystemInternalError),
		})
		c.Abort()
		return
	}
	now := time.Now()
	// session不存在
	if !b {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.NotLoginCode.Int(),
			Message: i18n.GetByKey(i18n.SystemNotLogin),
		})
		c.Abort()
		return
	}
	// 刷新token
	if session.ExpireAt < now.Add(apisession.RefreshSessionInterval).UnixMilli() {
		sessionStore.RefreshExpiry(sessionId, now.Add(apisession.SessionExpiry).UnixMilli())
	}
	c.Set(LoginUser, session.UserInfo)
	c.Next()
}

func GetLoginUser(c *gin.Context) (usermd.UserInfo, bool) {
	v, b := c.Get(LoginUser)
	if b {
		return v.(usermd.UserInfo), true
	}
	return usermd.UserInfo{}, false
}

func MustGetLoginUser(c *gin.Context) usermd.UserInfo {
	v := c.MustGet(LoginUser)
	return v.(usermd.UserInfo)
}

func GetSessionId(c *gin.Context) string {
	return c.GetHeader("Authorization")
}
