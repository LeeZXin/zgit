package projectapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/projectsrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/project", apicommon.CheckLogin)
		{
			group.POST("/insert", insertProject)
			group.POST("/list", listProject)
			group.POST("/delete", deleteProject)
			group.POST("/update", updateProject)
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
