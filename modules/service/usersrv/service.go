package usersrv

import (
	"context"
	"zgit/modules/model/usermd"
)

func GetUserInfoByPublicKey(ctx context.Context, pubKey string) (usermd.UserInfo, bool, error) {
	return usermd.UserInfo{
		Id:    "1",
		Name:  "zexin",
		Email: "zexin@fake.local",
	}, true, nil
}

func GetUserInfoByUserId(ctx context.Context, userId string) (usermd.UserInfo, bool, error) {
	return usermd.UserInfo{
		Id:    "1",
		Name:  "zexin",
		Email: "zexin@fake.local",
	}, true, nil
}
