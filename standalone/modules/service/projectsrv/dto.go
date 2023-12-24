package projectsrv

import (
	"zgit/standalone/modules/model/usermd"
	"zgit/util"
)

type InsertProjectReqDTO struct {
	Name     string
	Desc     string
	Operator usermd.UserInfo
}

func (r *InsertProjectReqDTO) IsValid() error {
	if len(r.Name) == 0 || len(r.Name) > 64 {
		return util.InvalidArgsError()
	}
	if len(r.Desc) > 128 {
		return util.InvalidArgsError()
	}
	if r.Operator.Account == "" {
		return util.InvalidArgsError()
	}
	return nil
}
