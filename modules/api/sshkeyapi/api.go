package sshkeyapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/modules/api/apicommon"
	"zgit/modules/model/sshkeymd"
	"zgit/modules/service/sshkeysrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sshKey", apicommon.CheckLogin)
		{
			group.POST("/delete", deleteSshKey)
			group.POST("/pubKey/insert", insertSshPubKey)
			group.POST("/pubKey/list")
			group.POST("/pubKey/verify")

			group.POST("/proxyKey/insert", insertSshProxyKey)
			group.POST("/proxyKey/list")
			group.POST("/proxyKey/verify")
		}
	})
}

func insertSshPubKey(c *gin.Context) {
	var req InsertSshPubKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.InsertSshKey(c.Request.Context(), sshkeysrv.InsertSshKeyReqDTO{
			Name:          req.Name,
			PubKeyContent: req.PubKeyContent,
			KeyType:       sshkeymd.UserPubKeyType,
			Operator:      apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}

func insertSshProxyKey(c *gin.Context) {
	var req InsertSshPubKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.InsertSshKey(c.Request.Context(), sshkeysrv.InsertSshKeyReqDTO{
			Name:          req.Name,
			PubKeyContent: req.PubKeyContent,
			KeyType:       sshkeymd.ProxyKeyType,
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
