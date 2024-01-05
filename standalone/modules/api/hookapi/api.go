package hookapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/pkg/apicode"
	"zgit/pkg/hook"
	"zgit/pkg/i18n"
	"zgit/setting"
	"zgit/standalone/modules/service/hooksrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		g := e.Group("/api/internal/hook", checkHookToken)
		{
			g.POST("/pre-receive", preReceive)
			g.POST("/post-receive", postReceive)
		}
	})
}

func checkHookToken(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization != setting.HookToken() {
		c.JSON(http.StatusUnauthorized, ginutil.BaseResp{
			Code:    apicode.UnauthorizedCode.Int(),
			Message: i18n.GetByKey(i18n.SystemUnauthorized),
		})
		c.Abort()
	} else {
		c.Next()
	}
}

func preReceive(c *gin.Context) {
	var req hook.Opts
	if util.ShouldBindJSON(&req, c) {
		logger.Logger.Info("pre-receive", req)
		err := hooksrv.PreReceive(c.Request.Context(), req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func postReceive(c *gin.Context) {
	var req hook.Opts
	if ginutil.ShouldBind(&req, c) {
		err := hooksrv.PostReceive(c.Request.Context(), req)
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
