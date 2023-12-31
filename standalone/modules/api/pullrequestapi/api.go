package pullrequestapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf-utils/timeutil"
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
			// 提交差异
			group.POST("/diffCommits", diffCommits)
			// 展示提交文件差异
			group.POST("/diffFile", diffFile)
			// 展示文件内容
			group.POST("/catFile", catFile)
			// 创建合并请求
			group.POST("/submit", submitPullRequest)
			// 关闭合并请求
			group.POST("/close", closePullRequest)
			// merge合并请求
			group.POST("/merge", mergePullRequest)
		}
	})
}

func diffFile(c *gin.Context) {
	var req DiffFileReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := pullrequestsrv.DiffFile(c.Request.Context(), pullrequestsrv.DiffFileReqDTO{
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
		ret.Lines, _ = listutil.Map(respDTO.Lines, func(t pullrequestsrv.DiffLineDTO) (DiffLineVO, error) {
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
		respDTO, err := pullrequestsrv.DiffCommits(c.Request.Context(), pullrequestsrv.DiffCommitsReqDTO{
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
		respVO.Commits, _ = listutil.Map(respDTO.Commits, func(t pullrequestsrv.CommitDTO) (CommitVO, error) {
			return commitDto2Vo(t), nil
		})
		respVO.DiffNumsStats.Stats, _ = listutil.Map(respDTO.DiffNumsStats.Stats, func(t pullrequestsrv.DiffNumsStatDTO) (DiffNumsStatVO, error) {
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

func commitDto2Vo(dto pullrequestsrv.CommitDTO) CommitVO {
	return CommitVO{
		Author:        dto.Author,
		Committer:     dto.Committer,
		AuthoredDate:  dto.AuthoredDate.Format(timeutil.DefaultTimeFormat),
		CommittedDate: dto.CommittedDate.Format(timeutil.DefaultTimeFormat),
		CommitMsg:     dto.CommitMsg,
		CommitId:      dto.CommitId,
		ShortId:       dto.ShortId,
	}
}

func catFile(c *gin.Context) {
	var req CatFileReqVO
	if util.ShouldBindJSON(&req, c) {
		lines, err := pullrequestsrv.CatFile(c.Request.Context(), pullrequestsrv.CatFileReqDTO{
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
		ret := CatFileRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		ret.Lines, _ = listutil.Map(lines, func(t pullrequestsrv.DiffLineDTO) (DiffLineVO, error) {
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
