package usermd

import (
	"context"
	"github.com/LeeZXin/zsf/xorm/xormutil"
)

func InsertUser(ctx context.Context, reqDTO InsertUserReqDTO) (User, error) {
	u := User{
		Account:   reqDTO.Account,
		Name:      reqDTO.Name,
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		AvatarUrl: reqDTO.AvatarUrl,
		IsAdmin:   reqDTO.IsAdmin,
	}
	_, err := xormutil.MustGetXormSession(ctx).Insert(&u)
	return u, err
}

func DeleteUser(ctx context.Context, user User) (bool, error) {
	rows, err := xormutil.MustGetXormSession(ctx).Where("account = ?", user.Account).Limit(1).Delete(new(User))
	return rows == 1, err
}

func GetByAccount(ctx context.Context, account string) (User, bool, error) {
	var ret User
	b, err := xormutil.MustGetXormSession(ctx).
		Where("account = ?", account).
		Get(&ret)
	return ret, b, err
}

func CountUser(ctx context.Context) (int64, error) {
	return xormutil.MustGetXormSession(ctx).Count(new(User))
}
