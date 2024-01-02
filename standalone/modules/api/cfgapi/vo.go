package cfgapi

import (
	"github.com/LeeZXin/zsf-utils/ginutil"
	"zgit/standalone/modules/service/cfgsrv"
)

type GetSysCfgRespVO struct {
	ginutil.BaseResp
	Cfg cfgsrv.SysCfg `json:"cfg"`
}

type UpdateSysCfgReqVO struct {
	cfgsrv.SysCfg
}
