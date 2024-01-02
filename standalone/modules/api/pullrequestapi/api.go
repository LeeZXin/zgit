package pullrequestapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/pullrequestsrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/pullRequest", apicommon.CheckLogin)
		{
			// 创建合并请求
			group.POST("/submit", submitPullRequest)
			// 关闭合并请求
			group.POST("/close", closePullRequest)
			// merge合并请求
			group.POST("/merge", mergePullRequest)
		}
	})
}

func submitPullRequest(c *gin.Context) {
	var req SubmitPullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.SubmitPullRequest(c.Request.Context(), pullrequestsrv.SubmitPullRequestReqDTO{
			RepoId:   req.RepoId,
			Target:   req.Target,
			Head:     req.Head,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func closePullRequest(c *gin.Context) {
	var req ClosePullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.ClosePullRequest(c.Request.Context(), pullrequestsrv.ClosePullRequestReqDTO{
			PrId:     req.PrId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func mergePullRequest(c *gin.Context) {
	var req MergePullRequestReqVO
	if util.ShouldBindJSON(&req, c) {
		err := pullrequestsrv.MergePullRequest(c.Request.Context(), pullrequestsrv.MergePullRequestReqDTO{
			PrId:     req.PrId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
