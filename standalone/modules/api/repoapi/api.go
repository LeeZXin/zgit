package repoapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
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
			// 获取仓库类型列表
			group.GET("/allTypeList", allTypeList)
			// 初始化仓库
			group.POST("/init", initRepo)
			// 删除仓库
			group.POST("/delete", deleteRepo)
			// 展示仓库列表
			group.POST("/list", listRepo)
			// 展示仓库主页
			group.POST("/tree", treeRepo)
			// 展示更多文件列表
			group.POST("/entries", entriesRepo)
			// 展示单个文件内容
			group.POST("/catFile", catFile)
			// 展示仓库所有分支
			group.POST("/allBranches", allBranches)
			// 展示仓库所有tag
			group.POST("/allTags", allTags)
			// gc
			group.POST("/gc", gc)
			// 提交差异
			group.POST("/diffCommits", diffCommits)
			// 展示提交文件差异
			group.POST("/diffFile", diffFile)
			// 展示文件内容
			group.POST("/showDiffTextContent", showDiffTextContent)
		}
		// 仓库管理
		group = e.Group("/api/repoManage", apicommon.CheckLogin)
		{
			group.POST("/insert")
			group.POST("/delete")
			group.POST("/list")
		}
	})
}

func allBranches(c *gin.Context) {
	var req AllBranchesReqVO
	if util.ShouldBindJSON(&req, c) {
		branches, err := reposrv.AllBranches(c.Request.Context(), reposrv.AllBranchesReqDTO{
			RepoId:   req.RepoId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, AllBranchesRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     branches,
		})
	}
}

func allTags(c *gin.Context) {
	var req AllTagsReqVO
	if util.ShouldBindJSON(&req, c) {
		branches, err := reposrv.AllTags(c.Request.Context(), reposrv.AllTagsReqDTO{
			RepoId:   req.RepoId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, AllBranchesRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     branches,
		})
	}
}

func gc(c *gin.Context) {
	var req GcReqVO
	if util.ShouldBindJSON(&req, c) {
		err := reposrv.Gc(c.Request.Context(), reposrv.GcReqDTO{
			RepoId:   req.RepoId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

// allTypeList 仓库类型列表
func allTypeList(c *gin.Context) {
	data, _ := listutil.Map(reposrv.AllTypeList(), func(t reposrv.RepoTypeDTO) (RepoTypeVO, error) {
		return RepoTypeVO{
			Option: t.Option,
			Name:   t.Name,
		}, nil
	})
	c.JSON(http.StatusOK, AllTypeListRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Data:     data,
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

// treeRepo 代码详情页
func treeRepo(c *gin.Context) {
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
				Offset:  repoRespDTO.Tree.Offset,
				Files:   fileDto2Vo(repoRespDTO.Tree.Files),
				Limit:   repoRespDTO.Tree.Limit,
				HasMore: repoRespDTO.Tree.HasMore,
			},
		})
	}
}

// entriesRepo 展示文件列表
func entriesRepo(c *gin.Context) {
	var req EntriesRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		repoRespDTO, err := reposrv.EntriesRepo(c.Request.Context(), reposrv.EntriesRepoReqDTO{
			RepoId:   req.RepoId,
			RefName:  req.RefName,
			Dir:      req.Dir,
			Offset:   req.Offset,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, TreeVO{
			Offset:  repoRespDTO.Offset,
			Files:   fileDto2Vo(repoRespDTO.Files),
			Limit:   repoRespDTO.Limit,
			HasMore: repoRespDTO.HasMore,
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

func listRepo(c *gin.Context) {
	var req ListRepoReqVO
	if util.ShouldBindJSON(&req, c) {
		repoList, err := reposrv.ListRepo(c.Request.Context(), reposrv.ListRepoReqDTO{
			ProjectId: req.ProjectId,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		repoVoList, _ := listutil.Map(repoList, func(t repomd.Repo) (RepoVO, error) {
			return RepoVO{
				Name:      t.Name,
				Path:      t.Path,
				Author:    t.Author,
				ProjectId: t.ProjectId,
				RepoType:  repomd.RepoType(t.RepoType).Readable(),
				IsEmpty:   t.IsEmpty,
				TotalSize: t.TotalSize,
				WikiSize:  t.WikiSize,
				GitSize:   t.GitSize,
				LfsSize:   t.LfsSize,
				Created:   t.Created.Format(timeutil.DefaultTimeFormat),
			}, nil
		})
		c.JSON(http.StatusOK, ListRepoRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			RepoList: repoVoList,
		})
	}
}

func catFile(c *gin.Context) {
	var req CatFileReqVO
	if util.ShouldBindJSON(&req, c) {
		fileMode, content, err := reposrv.CatFile(c.Request.Context(), reposrv.CatFileReqDTO{
			RepoId:   req.RepoId,
			RefName:  req.RefName,
			Dir:      req.Dir,
			FileName: req.FileName,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, CatFileRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Mode:     fileMode.Readable(),
			Content:  content,
		})
	}
}

func diffFile(c *gin.Context) {
	var req DiffFileReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.DiffFile(c.Request.Context(), reposrv.DiffFileReqDTO{
			RepoId:   req.RepoId,
			Target:   req.Target,
			Head:     req.Head,
			FileName: req.FileName,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := DiffFileRespVO{
			FilePath:    respDTO.FilePath,
			OldMode:     respDTO.OldMode,
			Mode:        respDTO.Mode,
			IsSubModule: respDTO.IsSubModule,
			FileType:    respDTO.FileType.String(),
			IsBinary:    respDTO.IsBinary,
			RenameFrom:  respDTO.RenameFrom,
			RenameTo:    respDTO.RenameTo,
			CopyFrom:    respDTO.CopyFrom,
			CopyTo:      respDTO.CopyTo,
		}
		ret.Lines, _ = listutil.Map(respDTO.Lines, func(t reposrv.DiffLineDTO) (DiffLineVO, error) {
			return DiffLineVO{
				Index:   t.Index,
				LeftNo:  t.LeftNo,
				Prefix:  t.Prefix,
				RightNo: t.RightNo,
				Text:    t.Text,
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func diffCommits(c *gin.Context) {
	var req PrepareMergeReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := reposrv.DiffCommits(c.Request.Context(), reposrv.DiffCommitsReqDTO{
			RepoId:   req.RepoId,
			Target:   req.Target,
			Head:     req.Head,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := PrepareMergeRespVO{
			BaseResp:     ginutil.DefaultSuccessResp,
			Target:       respDTO.Target,
			Head:         respDTO.Head,
			TargetCommit: commitDto2Vo(respDTO.TargetCommit),
			HeadCommit:   commitDto2Vo(respDTO.HeadCommit),
			NumFiles:     respDTO.NumFiles,
			DiffNumsStats: DiffNumsStatInfoVO{
				FileChangeNums: respDTO.DiffNumsStats.FileChangeNums,
				InsertNums:     respDTO.DiffNumsStats.InsertNums,
				DeleteNums:     respDTO.DiffNumsStats.DeleteNums,
			},
			ConflictFiles: respDTO.ConflictFiles,
			CanMerge:      respDTO.CanMerge,
		}
		respVO.Commits, _ = listutil.Map(respDTO.Commits, func(t reposrv.CommitDTO) (CommitVO, error) {
			return commitDto2Vo(t), nil
		})
		respVO.DiffNumsStats.Stats, _ = listutil.Map(respDTO.DiffNumsStats.Stats, func(t reposrv.DiffNumsStatDTO) (DiffNumsStatVO, error) {
			return DiffNumsStatVO{
				RawPath:    t.RawPath,
				Path:       t.Path,
				TotalNums:  t.TotalNums,
				InsertNums: t.InsertNums,
				DeleteNums: t.DeleteNums,
			}, nil
		})
		c.JSON(http.StatusOK, respVO)
	}
}

func showDiffTextContent(c *gin.Context) {
	var req ShowDiffTextContentReqVO
	if util.ShouldBindJSON(&req, c) {
		lines, err := reposrv.ShowDiffTextContent(c.Request.Context(), reposrv.ShowDiffTextContentReqDTO{
			RepoId:    req.RepoId,
			CommitId:  req.CommitId,
			FileName:  req.FileName,
			Offset:    req.Offset,
			Limit:     req.Limit,
			Direction: req.Direction,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ShowDiffTextContentRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		ret.Lines, _ = listutil.Map(lines, func(t reposrv.DiffLineDTO) (DiffLineVO, error) {
			return DiffLineVO{
				Index:   t.Index,
				LeftNo:  t.LeftNo,
				Prefix:  t.Prefix,
				RightNo: t.RightNo,
				Text:    t.Text,
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}
