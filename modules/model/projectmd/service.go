package projectmd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
	"strconv"
)

func GetByProjectId(ctx context.Context, projectId string) (Project, bool, error) {
	var ret Project
	b, err := mysqlstore.GetXormSession(ctx).Where("project_id = ?", projectId).Get(&ret)
	return ret, b, err
}

func ListProjectByCorpId(ctx context.Context, reqDTO ListProjectByCorpIdReqDTO) (ListProjectByCorpIdRespDTO, error) {
	ret := make([]Project, 0)
	session := mysqlstore.GetXormSession(ctx).Where("corp_id = ?", reqDTO.CorpId)
	if reqDTO.Cursor != "" {
		cursor, _ := strconv.ParseInt(reqDTO.Cursor, 10, 64)
		if cursor > 0 && reqDTO.Limit > 0 {
			session.And("id > ?", cursor).Limit(reqDTO.Limit)
		}
	}
	err := session.Find(&ret)
	if err != nil {
		return ListProjectByCorpIdRespDTO{}, err
	}
	respDTO := ListProjectByCorpIdRespDTO{
		Data: ret,
	}
	if len(ret) > 0 {
		respDTO.Next = strconv.FormatInt(ret[len(ret)-1].Id, 10)
	}
	return respDTO, nil
}
