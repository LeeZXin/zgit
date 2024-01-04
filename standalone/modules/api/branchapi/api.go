package branchapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/branchsrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		// 保护分支
		group := e.Group("/api/protectedBranch", apicommon.CheckLogin)
		{
			// 新增保护分支
			group.POST("/insert", insertProtectedBranch)
			group.POST("/delete", deleteProtectedBranch)
			group.POST("/list", listProtectedBranch)
		}
	})
}

func insertProtectedBranch(c *gin.Context) {
	var req InsertProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.InsertProtectedBranch(c.Request.Context(), branchsrv.InsertProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Branch:   req.Branch,
			Cfg:      req.Cfg,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func deleteProtectedBranch(c *gin.Context) {
	var req DeleteProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		err := branchsrv.DeleteProtectedBranch(c.Request.Context(), branchsrv.DeleteProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Branch:   req.Branch,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func listProtectedBranch(c *gin.Context) {
	var req ListProtectedBranchReqVO
	if util.ShouldBindJSON(&req, c) {
		branchList, err := branchsrv.ListProtectedBranch(c.Request.Context(), branchsrv.ListProtectedBranchReqDTO{
			RepoId:   req.RepoId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		respVO := ListProtectedBranchRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
		}
		respVO.Branches, _ = listutil.Map(branchList, func(t branchsrv.ProtectedBranchDTO) (ProtectedBranchVO, error) {
			return ProtectedBranchVO{
				Branch: t.Branch,
				Cfg:    t.Cfg,
			}, nil
		})
		c.JSON(http.StatusOK, respVO)
	}
}
