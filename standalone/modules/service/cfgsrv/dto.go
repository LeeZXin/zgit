package cfgsrv

import (
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type UpdateSysCfgReqDTO struct {
	SysCfg
	Operator usermd.UserInfo
}

func (r *UpdateSysCfgReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}

type GetSysCfgReqDTO struct {
	Operator usermd.UserInfo
}

func (r *GetSysCfgReqDTO) IsValid() error {
	if !util.ValidateOperator(r.Operator) {
		return util.InvalidArgsError()
	}
	return nil
}
