package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/model/repomd"
	"zgit/standalone/modules/service/reposrv"
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
			// 展示仓库主页
			group.GET("/tree", tree)

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

func tree(c *gin.Context) {
	var req TreeRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		repoRespDTO, err := reposrv.TreeRepo(c.Request.Context(), reposrv.TreeRepoReqDTO{
			RepoId:   req.RepoId,
			RefName:  req.RefName,
			Dir:      req.Dir,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, TreeRepoRespVO{
			BaseResp:     ginutil.DefaultSuccessResp,
			IsEmpty:      repoRespDTO.IsEmpty,
			ReadmeText:   repoRespDTO.ReadmeText,
			RecentCommit: commitDto2Vo(repoRespDTO.RecentCommit),
			Tree: TreeVO{
				Files:   fileDto2Vo(repoRespDTO.Tree.Files),
				Limit:   repoRespDTO.Tree.Limit,
				HasMore: repoRespDTO.Tree.HasMore,
			},
		})
	}
}

func commitDto2Vo(dto reposrv.CommitDTO) CommitVO {
	return CommitVO{
		Author:        dto.Author,
		Committer:     dto.Committer,
		AuthoredDate:  util.ReadableTimeComparingNow(dto.AuthoredDate),
		CommittedDate: util.ReadableTimeComparingNow(dto.CommittedDate),
		CommitMsg:     dto.CommitMsg,
		CommitId:      dto.CommitId,
		ShortId:       dto.ShortId,
	}
}

func fileDto2Vo(list []reposrv.FileDTO) []FileVO {
	ret, _ := listutil.Map(list, func(t reposrv.FileDTO) (FileVO, error) {
		return FileVO{
			Mode:    t.Mode,
			RawPath: t.RawPath,
			Path:    t.Path,
			Commit:  commitDto2Vo(t.Commit),
		}, nil
	})
	return ret
}

func initRepo(c *gin.Context) {
	var req InitRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.InitRepo(c.Request.Context(), reposrv.InitRepoReqDTO{
			Operator:      apicommon.MustGetLoginUser(c),
			Name:          req.Name,
			Desc:          req.Desc,
			RepoType:      repomd.RepoType(req.RepoType),
			CreateReadme:  req.CreateReadme,
			GitIgnoreName: req.GitIgnoreName,
			DefaultBranch: req.DefaultBranch,
			ProjectId:     req.ProjectId,
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
