package sshkeyapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf-utils/listutil"
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
			// 删除
			group.POST("/delete", deleteSshKey)
			// 插入
			group.POST("/insert", insertSshKey)
			// 列表展示
			group.POST("/list", listSshKey)
			// 校验
			group.POST("/verify", verifySshKey)
			// 获取校验token
			group.POST("/getToken", getToken)
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

func listSshKey(c *gin.Context) {
	var req ListSshKeyReqVO
	if util.ShouldBindJSON(&req, c) {
		respDTO, err := sshkeysrv.ListSshKey(c.Request.Context(), sshkeysrv.ListSshKeyReqDTO{
			Offset:   req.Offset,
			Limit:    req.Limit,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		ret := ListSshKeyRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Cursor:   respDTO.Cursor,
		}
		ret.Data, _ = listutil.Map(respDTO.KeyList, func(t sshkeysrv.SshKeyDTO) (SshKeyVO, error) {
			return SshKeyVO{
				KeyId:       t.KeyId,
				Name:        t.Name,
				Fingerprint: t.Fingerprint,
			}, nil
		})
		c.JSON(http.StatusOK, ret)
	}
}

func getToken(c *gin.Context) {
	var req GetTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		token, err := sshkeysrv.GetToken(c.Request.Context(), sshkeysrv.GetTokenReqDTO{
			KeyId:    req.KeyId,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, GetTokenRespVO{
			BaseResp: ginutil.DefaultSuccessResp,
			Token:    token,
		})
	}
}

// verifySshKey 校验sshKey
func verifySshKey(c *gin.Context) {
	var req VerifyTokenReqVO
	if util.ShouldBindJSON(&req, c) {
		err := sshkeysrv.VerifySshKey(c.Request.Context(), sshkeysrv.VerifySshKeyReqDTO{
			KeyId:     req.KeyId,
			Token:     req.Token,
			Signature: req.Signature,
			Operator:  apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
