package cfgapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"github.com/LeeZXin/zsf/http/httpserver"
	"github.com/gin-gonic/gin"
	"net/http"
	"zgit/standalone/modules/api/apicommon"
	"zgit/standalone/modules/service/cfgsrv"
	"zgit/util"
)

func InitApi() {
	httpserver.AppendRegisterRouterFunc(func(e *gin.Engine) {
		group := e.Group("/api/sysCfg", apicommon.CheckLogin)
		{
			group.GET("/get", getSysCfg)
			group.POST("/update", updateSysCfg)
		}
	})
}

func getSysCfg(c *gin.Context) {
	cfg, err := cfgsrv.GetSysCfg(c.Request.Context(), cfgsrv.GetSysCfgReqDTO{
		Operator: apicommon.MustGetLoginUser(c),
	})
	if err != nil {
		util.HandleApiErr(err, c)
		return
	}
	c.JSON(http.StatusOK, GetSysCfgRespVO{
		BaseResp: ginutil.DefaultSuccessResp,
		Cfg:      cfg,
	})
}

func updateSysCfg(c *gin.Context) {
	var req UpdateSysCfgReqVO
	if util.ShouldBindJSON(&req, c) {
		err := cfgsrv.UpdateSysCfg(c.Request.Context(), cfgsrv.UpdateSysCfgReqDTO{
			SysCfg:   req.SysCfg,
			Operator: apicommon.MustGetLoginUser(c),
		})
		if err != nil {
			util.HandleApiErr(err, c)
			return
		}
		c.JSON(http.StatusOK, ginutil.DefaultSuccessResp)
	}
}
