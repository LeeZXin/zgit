package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/modules/api/apicommon"
	"zgit/modules/model/repomd"
	"zgit/modules/service/reposrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/repo", apicommon.CheckLogin)
		{
			// 获取模版列表
			group.GET("/allGitIgnoreTemplateList", allGitIgnoreTemplateList)
			// 初始化仓库
			group.POST("/init", initRepo)
			// 删除仓库
			group.POST("/delete", deleteRepo)
			// 展示仓库列表
			group.POST("/list")
		}
	})
}

// allGitIgnoreTemplateList 获取模版列表
func allGitIgnoreTemplateList(c *gin.Context) {
	templateList := reposrv.AllGitIgnoreTemplateList()
	c.JSON(http.StatusOK, AllGitIgnoreTemplateListRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     templateList,
	})
}

func initRepo(c *gin.Context) {
	var req InitRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.InitRepo(c.Request.Context(), reposrv.InitRepoReqDTO{
			Operator:      apicommon.MustGetLoginUser(c),
			RepoName:      req.RepoName,
			RepoDesc:      req.RepoDesc,
			RepoType:      repomd.RepoType(req.RepoType),
			CreateReadme:  req.CreateReadme,
			GitIgnoreName: req.GitIgnoreName,
			DefaultBranch: req.DefaultBranch,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func deleteRepo(c *gin.Context) {
	var req DeleteRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.DeleteRepo(c.Request.Context(), reposrv.DeleteRepoReqDTO{
			Operator: apicommon.MustGetLoginUser(c),
			RepoId:   req.RepoId,
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
