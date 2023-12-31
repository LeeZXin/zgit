package userapi

import "github.com/LeeZXin/zsf-utils/ginutil"

type LoginReqVO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginRespVO struct {
	ginutil.BaseResp
	SessionId string `json:"sessionId"`
}

type InsertUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsAdmin   bool   `json:"isAdmin"`
	AvatarUrl string `json:"avatarUrl"`
}

type RegisterUserReqVO struct {
	Account   string `json:"account"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatarUrl"`
}

type DeleteUserReqVO struct {
	Account string `json:"account"`
}

type UpdateUserReqVO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"isAdmin"`
	AvatarUrl string `json:"avatarUrl"`
}
