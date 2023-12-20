package hookapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		g := e.Group("/api/internal/hook")
		{
			g.POST("/pre-receive", func(c *gin.Context) {
				var req OptsReqVO
				if ginutil.ShouldBind(&req, c) {
					logger.Logger.Info("pre-receive", req)
					c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
				}
			})
			g.POST("/post-receive", func(c *gin.Context) {
				var req OptsReqVO
				if ginutil.ShouldBind(&req, c) {
					logger.Logger.Info("post-receive", req)
					c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
				}
			})
		}
	})
}
