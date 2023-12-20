package usermd

import (
	"context"
	"github.com/LeeZXin/zsf-utils/idutil"
	"github.com/LeeZXin/zsf/xorm/mysqlstore"
)

func genUserId() string {
	return idutil.RandomUuid()
}

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) (User, error) {
	u := User{
		UserId:    genUserId(),
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		CorpId:    reqDTO.CorpId,
		AvatarUrl: reqDTO.AvatarUrl,
		IsAdmin:   reqDTO.IsAdmin,
	}
	_, err := mysqlstore.GetXormSession(ctx).Insert(&u)
	return u, err
}

func DeleteUser(ctx context.Context, user User) (bool, error) {
	rows, err := mysqlstore.GetXormSession(ctx).Where("user_id = ?", user.UserId).Limit(1).Delete(new(User))
	return rows == 1, err
}

func GetByAccountAndCorpId(ctx context.Context, account string, corpId string) (User, bool, error) {
	var ret User
	b, err := mysqlstore.GetXormSession(ctx).
		Where("account = ?", account).
		And("corp_id = ?", corpId).
		Get(&ret)
	return ret, b, err
}

func GetByUserId(ctx context.Context, userId string) (User, bool, error) {
	var ret User
	b, err := mysqlstore.GetXormSession(ctx).
		Where("user_id = ?", userId).
		Get(&ret)
	return ret, b, err
}
