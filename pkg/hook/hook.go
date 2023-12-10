package hook

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ApiPreReceiveUrl  = "api/internal/hook/pre-receive"
	ApiPostReceiveUrl = "api/internal/hook/post-receive"
)

type RevInfo struct {
	OldCommitId string `json:"oldCommitId"`
	NewCommitId string `json:"newCommitId"`
	RefName     string `json:"refName"`
}

type OptsReqVO struct {
	RevInfoList []RevInfo `json:"revInfoList"`
	RepoId      string    `json:"repoId"`
	PrId        string    `json:"prId"`
	PusherId    string    `json:"pusherId"`
	IsWiki      bool      `json:"isWiki"`
}

type HttpRespVO struct {
	ginutil.BaseResp
}

func init() {
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
