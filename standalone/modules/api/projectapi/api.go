package projectapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/projectsrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		// 项目
		group := e.Group("/api/project", apicommon.CheckLogin)
		{
			group.POST("/insert", insertProject)
			group.POST("/list", listProject)
			group.POST("/delete", deleteProject)
			group.POST("/update", updateProject)
		}
		// 项目用户
		group = e.Group("/api/projectUser", apicommon.CheckLogin)
		{
			group.POST("/upsert", upsertProjectUser)
			group.POST("/list", listProjectUser)
			group.POST("/delete", deleteProjectUser)
		}
		// 项目用户组
		group = e.Group("/api/projectUserGroup", apicommon.CheckLogin)
		{
			group.POST("/insert", insertProjectUserGroup)
			group.POST("/list", listProjectUserGroup)
			group.POST("/updateName", updateProjectUserGroupName)
			group.POST("/updatePerm", updateProjectUserGroupPerm)
			group.POST("/delete", deleteProjectUserGroup)
		}
	})
}

func insertProject(c *gin.Context) {
	var req InsertProjectReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.InsertProject(c.Request.Context(), projectsrv.InsertProjectReqDTO{
			Name:     req.Name,
			Desc:     req.Desc,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func listProject(c *gin.Context) {

}

func deleteProject(c *gin.Context) {

}

func updateProject(c *gin.Context) {

}

func upsertProjectUser(c *gin.Context) {
	var req UpsertProjectUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.UpsertProjectUser(c.Request.Context(), projectsrv.UpsertProjectUserReqDTO{
			ProjectId: req.ProjectId,
			Account:   req.Account,
			GroupId:   req.GroupId,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func listProjectUser(c *gin.Context) {

}

func deleteProjectUser(c *gin.Context) {
	var req DeleteProjectUserReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.DeleteProjectUser(c.Request.Context(), projectsrv.DeleteProjectUserReqDTO{
			ProjectId: req.ProjectId,
			Account:   req.Account,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func insertProjectUserGroup(c *gin.Context) {
	var req InsertProjectUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.InsertProjectUserGroup(c.Request.Context(), projectsrv.InsertProjectUserGroupReqDTO{
			ProjectId: req.ProjectId,
			Name:      req.Name,
			Perm:      req.Perm,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func updateProjectUserGroupName(c *gin.Context) {
	var req UpdateProjectUserGroupNameReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.UpdateProjectUserGroupName(c.Request.Context(), projectsrv.UpdateProjectUserGroupNameReqDTO{
			GroupId:  req.GroupId,
			Name:     req.Name,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func updateProjectUserGroupPerm(c *gin.Context) {
	var req UpdateProjectUserGroupPermReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.UpdateProjectUserGroupPerm(c.Request.Context(), projectsrv.UpdateProjectUserGroupPermReqDTO{
			GroupId:  req.GroupId,
			Perm:     req.Perm,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func deleteProjectUserGroup(c *gin.Context) {
	var req DeleteProjectUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		err := projectsrv.DeleteProjectUserGroup(c.Request.Context(), projectsrv.DeleteProjectUserGroupReqDTO{
			GroupId:  req.GroupId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func listProjectUserGroup(c *gin.Context) {
	var req ListProjectUserGroupReqVO
	if util.ShouldBindJSON(&req, c) {
		groups, err := projectsrv.ListProjectUserGroup(c.Request.Context(), projectsrv.ListProjectUserGroupReqDTO{
			ProjectId: req.ProjectId,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret, _ := listutil.Map(groups, func(t projectsrv.ProjectUserGroupDTO) (ProjectUserGroupVO, error) {
			return ProjectUserGroupVO{
				GroupId:   t.GroupId,
				ProjectId: t.ProjectId,
				Name:      t.Name,
				Perm:      t.Perm,
			}, nil
		})
		c.JSON(http.StatusOK, ListProjectUserGroupRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Data:     ret,
		})
	}
}
