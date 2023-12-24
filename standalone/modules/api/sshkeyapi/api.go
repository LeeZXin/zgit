package sshkeyapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/sshkeysrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sshKey", apicommon.CheckLogin)
		{
			group.POST("/delete", deleteSshKey)
			group.POST("/insert", insertSshKey)
			group.POST("/list")
			group.POST("/verify")
		}
	})
}

func insertSshKey(c *gin.Context) {
	var req InsertSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.InsertSshKey(c.Request.Context(), sshkeysrv.InsertSshKeyReqDTO{
			Name:          req.Name,
			PubKeyContent: req.PubKeyContent,
			Operator:      apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func deleteSshKey(c *gin.Context) {
	var req DeleteSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.DeleteSshKey(c.Request.Context(), sshkeysrv.DeleteSshKeyReqDTO{
			KeyId:    req.KeyId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
